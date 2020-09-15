// Package setup provides common logic for configuring the various services.
package setup

import (
	"context"
	"fmt"

	"github.com/alienvspredator/simple-tgbot/internal/logging"
	"github.com/alienvspredator/simple-tgbot/internal/secrets"
	"github.com/alienvspredator/simple-tgbot/internal/serverenv"
	"github.com/sethvargo/go-envconfig"
)

// SecretManagerConfigProvider signals that the config knows how to configure a
// secret manager.
type SecretManagerConfigProvider interface {
	SecretManagerConfig() *secrets.Config
}

// Setup runs common initialization code for all servers. See SetupWith.
func Setup(ctx context.Context, config interface{}) (*serverenv.ServerEnv, error) {
	return SetupWith(ctx, config, envconfig.OsLookuper())
}

// SetupWith processes the given configuration using envconfig. It is
// responsible for establishing database connections, resolving secrets, and
// accessing app configs. The provided interface must implement the various
// interfaces.
func SetupWith(ctx context.Context, config interface{}, l envconfig.Lookuper) (
	*serverenv.ServerEnv, error,
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
			return nil, fmt.Errorf("unable to process secret manager env: %w", err)
		}

		var err error
		sm, err = secrets.SecretManagerFor(ctx, smConfig.SecretManagerType)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to secret manager: %w", err)
		}

		// Enable caching, if a TTL was provided.
		if ttl := smConfig.SecretCacheTTL; ttl > 0 {
			sm, err = secrets.WrapCacher(ctx, sm, ttl)
			if err != nil {
				return nil, fmt.Errorf("unable to create secret manager cache: %w", err)
			}
		}

		// Enable secret expansion, if enabled.
		if smConfig.SecretExpansion {
			sm, err = secrets.WrapJSONExpander(ctx, sm)
			if err != nil {
				return nil, fmt.Errorf("unable to create json expander secret manager: %w", err)
			}
		}

		// Update the mutators to process secrets.
		mutatorFuncs = append(mutatorFuncs, secrets.Resolver(sm, smConfig))

		// Update serverEnv setup.
		serverEnvOpts = append(serverEnvOpts, serverenv.WithSecretManager(sm))

		log.Infow("secret manager", "config", smConfig)
	}

	if err := envconfig.ProcessWith(ctx, config, l, mutatorFuncs...); err != nil {
		return nil, fmt.Errorf("load environment vars: %w", err)
	}
	log.Infow("provided", "config", config)

	return serverenv.New(ctx, serverEnvOpts...), nil
}
