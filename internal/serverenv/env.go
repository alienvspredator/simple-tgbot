// Package serverenv defines common parameters for the server environment.
package serverenv

import (
	"context"

	"github.com/alienvspredator/simple-tgbot/internal/database"
	"github.com/alienvspredator/simple-tgbot/internal/observability"
	"github.com/alienvspredator/simple-tgbot/internal/secrets"
	"github.com/hashicorp/go-multierror"
)

// ServerEnv represents environment configuration in this application.
type ServerEnv struct {
	database      *database.DB
	secretManager secrets.SecretManager
	// oe observability exporter
	oe observability.Exporter
}

// Option defines function types to modify the ServerEnv on creation.
type Option func(env *ServerEnv) *ServerEnv

// WithSecretManager creates an Option to install a specific secret manager to use.
func WithSecretManager(sm secrets.SecretManager) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.secretManager = sm
		return s
	}
}

func WithObservabilityExporter(oe observability.Exporter) Option {
	return func(env *ServerEnv) *ServerEnv {
		env.oe = oe
		return env
	}
}

// WithDatabase attached a database to the environment.
func WithDatabase(db *database.DB) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.database = db
		return s
	}
}

// New creates a new ServerEnv with the requested options.
func New(_ context.Context, opts ...Option) *ServerEnv {
	env := &ServerEnv{}

	for _, f := range opts {
		env = f(env)
	}

	return env
}

func (s *ServerEnv) SecretManager() secrets.SecretManager {
	return s.secretManager
}

func (s *ServerEnv) Database() *database.DB {
	return s.database
}

func (s *ServerEnv) ObservabilityExporter() observability.Exporter {
	return s.oe
}

// Close shuts down the server env, closing database connections, etc.
func (s *ServerEnv) Close(ctx context.Context) (err error) {
	if s == nil {
		return nil
	}

	if s.database != nil {
		s.database.Close(ctx)
	}

	if s.oe != nil {
		if closeErr := s.oe.Close(); closeErr != nil {
			err = multierror.Append(err, closeErr)
		}
	}

	return
}
