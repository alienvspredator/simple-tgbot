// Package serverenv defines common parameters for the server environment.
package serverenv

import "context"

// ServerEnv represents environment configuration in this application.
type ServerEnv struct {
}

// Option defines function types to modify the ServerEnv on creation.
type Option func(env *ServerEnv) *ServerEnv

// New creates a new ServerEnv with the requested options.
func New(ctx context.Context, opts ...Option) *ServerEnv {
	env := &ServerEnv{}

	for _, f := range opts {
		env = f(env)
	}

	return env
}

// Close shuts down the server env, closing database connections, etc.
func (s *ServerEnv) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}

	return nil
}
