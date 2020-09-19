// Package serverenv defines common parameters for the server environment.
package serverenv

import (
	"context"

	"github.com/alienvspredator/simple-tgbot/internal/database"
	"github.com/alienvspredator/simple-tgbot/internal/secrets"
)

// ServerEnv represents environment configuration in this application.
type ServerEnv struct {
	database      *database.DB
	secretManager secrets.SecretManager
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

// WithDatabase attached a database to the environment.
func WithDatabase(db *database.DB) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.database = db
		return s
	}
}

// New creates a new ServerEnv with the requested options.
func New(ctx context.Context, opts ...Option) *ServerEnv {
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

// Close shuts down the server env, closing database connections, etc.
func (s *ServerEnv) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}

	if s.database != nil {
		s.database.Close(ctx)
	}

	return nil
}
