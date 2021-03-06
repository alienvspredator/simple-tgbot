// Package setup provides common logic for configuring the various services.
package setup

import (
	"context"
	"fmt"

	"github.com/alienvspredator/simple-tgbot/internal/database"
	"github.com/alienvspredator/simple-tgbot/internal/logging"
	"github.com/alienvspredator/simple-tgbot/internal/observability"
	"github.com/alienvspredator/simple-tgbot/internal/observability/views"
	"github.com/alienvspredator/simple-tgbot/internal/secrets"
	"github.com/alienvspredator/simple-tgbot/internal/serverenv"
	"github.com/alienvspredator/simple-tgbot/internal/zapfluentd"
	"github.com/sethvargo/go-envconfig"
)

// DatabaseConfigProvider ensures that the environment config can provide a DB config.
// All binaries in this application connect to the database via the same method.
type DatabaseConfigProvider interface {
	DatabaseConfig() *database.Config
}

// ObservabilityExporterConfigProvider signals that the config knows how to configure an
// observability exporter.
type ObservabilityExporterConfigProvider interface {
	ObservabilityExporterConfig() *observability.Config
}

// SecretManagerConfigProvider signals that the config knows how to configure a
// secret manager.
type SecretManagerConfigProvider interface {
	SecretManagerConfig() *secrets.Config
}

type FluentConfigProvider interface {
	FluentConfig() *zapfluentd.Config
}

// Setup runs common initialization code for all servers. See SetupWith.
func Setup(ctx context.Context, config interface{}) (context.Context, *serverenv.ServerEnv, error) {
	return SetupWith(ctx, config, envconfig.OsLookuper())
}

// SetupWith processes the given configuration using envconfig. It is
// responsible for establishing database connections, resolving secrets, and
// accessing app configs. The provided interface must implement the various
// interfaces.
func SetupWith(ctx context.Context, config interface{}, l envconfig.Lookuper) (
	context.Context, *serverenv.ServerEnv, error,
) {
	log := logging.FromContext(ctx)

	var mutatorFuncs []envconfig.MutatorFunc

	var serverEnvOpts []serverenv.Option

	// Load the secret manager - this needs to be loaded first because other
	// processors may require access to secrets.
	var sm secrets.SecretManager
	if provider, ok := config.(SecretManagerConfigProvider); ok {
		log.Info("configuring secret manager")

		// The environment configuration defines which secret manager to use, so we
		// need to process the envconfig in here. Once we know which secret manager
		// to use, we can load the secret manager and then re-process (later) any
		// secret:// references.
		//
		// NOTE: it is not possible to specify which secret manager to use via a
		// secret:// reference. This configuration option must always be the
		// plaintext string.
		smConfig := provider.SecretManagerConfig()
		if err := envconfig.ProcessWith(ctx, smConfig, l, mutatorFuncs...); err != nil {
			return ctx, nil, fmt.Errorf("unable to process secret manager env: %w", err)
		}

		var err error
		sm, err = secrets.SecretManagerFor(ctx, smConfig.SecretManagerType)
		if err != nil {
			return ctx, nil, fmt.Errorf("unable to connect to secret manager: %w", err)
		}

		// Enable caching, if a TTL was provided.
		if ttl := smConfig.SecretCacheTTL; ttl > 0 {
			sm, err = secrets.WrapCacher(ctx, sm, ttl)
			if err != nil {
				return ctx, nil, fmt.Errorf("unable to create secret manager cache: %w", err)
			}
		}

		// Enable secret expansion, if enabled.
		if smConfig.SecretExpansion {
			sm, err = secrets.WrapJSONExpander(ctx, sm)
			if err != nil {
				return ctx, nil, fmt.Errorf(
					"unable to create json expander secret manager: %w", err,
				)
			}
		}

		// Update the mutators to process secrets.
		mutatorFuncs = append(mutatorFuncs, secrets.Resolver(sm, smConfig))

		// Update serverEnv setup.
		serverEnvOpts = append(serverEnvOpts, serverenv.WithSecretManager(sm))

		log.Infow("secret manager", "config", smConfig)
	}

	if err := envconfig.ProcessWith(ctx, config, l, mutatorFuncs...); err != nil {
		return ctx, nil, fmt.Errorf("load environment vars: %w", err)
	}
	log.Infow("provided", "config", config)

	if provider, ok := config.(FluentConfigProvider); ok {
		log.Info("configuring fluent")

		fConf := provider.FluentConfig()
		if fConf.Host != "" {
			ws, err := zapfluentd.NewWriteSyncer(fConf)
			if err != nil {
				return ctx, nil, fmt.Errorf("zap fluentd WriteSyncer initializing: %w", err)
			}
			ctx = logging.WithLogger(ctx, logging.NewLogger(logging.WithWriteSyncer(ws)))
		}
	}

	if provider, ok := config.(ObservabilityExporterConfigProvider); ok {
		log.Info("configuring observability exporter")

		oeConfig := provider.ObservabilityExporterConfig()
		oe, err := observability.NewFromEnv(ctx, oeConfig)
		if err != nil {
			return ctx, nil, fmt.Errorf("create observability provider: %w", err)
		}
		if err := oe.StartExporter(views.Register); err != nil {
			return nil, nil, fmt.Errorf("start observability: %w", err)
		}
		serverEnvOpts = append(serverEnvOpts, serverenv.WithObservabilityExporter(oe))

		log.Infow("observability", "config", oeConfig)
	}

	// Setup the database connection
	if provider, ok := config.(DatabaseConfigProvider); ok {
		log.Info("configuring database")

		dbConfig := provider.DatabaseConfig()
		db, err := database.NewFromEnv(ctx, dbConfig)
		if err != nil {
			return ctx, nil, fmt.Errorf("database connecting: %w", err)
		}

		serverEnvOpts = append(serverEnvOpts, serverenv.WithDatabase(db))
		log.Infow("database", "config", dbConfig)
	}

	return ctx, serverenv.New(ctx, serverEnvOpts...), nil
}
