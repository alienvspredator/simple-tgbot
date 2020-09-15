// Package setup provides common logic for configuring the various services.
package setup

import (
	"context"
	"fmt"

	"github.com/alienvspredator/simple-tgbot/internal/logging"
	"github.com/alienvspredator/simple-tgbot/internal/serverenv"
	"github.com/sethvargo/go-envconfig"
)

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

	if err := envconfig.ProcessWith(ctx, config, l, mutatorFuncs...); err != nil {
		return nil, fmt.Errorf("load environment vars: %w", err)
	}
	log.Infow("provided", "config", config)

	return serverenv.New(ctx, serverEnvOpts...), nil
}
