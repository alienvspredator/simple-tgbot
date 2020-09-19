package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alienvspredator/simple-tgbot/internal/logging"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

// NewFromEnv sets up the database connections using the configuration in the
// process's environment variables. This should be called just once per server
// instance.
func NewFromEnv(ctx context.Context, config *Config) (*DB, error) {
	log := logging.FromContext(ctx)
	log.Info("creating db connection pool")

	connStr := dbConnectionString(config)

	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close releases database connections.
func (db *DB) Close(ctx context.Context) {
	log := logging.FromContext(ctx)
	log.Info("closing db connection pool")
	db.Pool.Close()
}

// dbConnectionString builds a connection string suitable for the pgx Postgres driver, using the
// values of vars.
func dbConnectionString(config *Config) string {
	vals := dbValues(config)
	var p []string
	for k, v := range vals {
		p = append(p, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(p, " ")
}

func setIfNotEmpty(m map[string]string, key, val string) {
	if val != "" {
		m[key] = val
	}
}

func setIfPositive(m map[string]string, key string, val int) {
	if val > 0 {
		m[key] = fmt.Sprintf("%d", val)
	}
}

func setIfPositiveDuration(m map[string]string, key string, d time.Duration) {
	if d > 0 {
		m[key] = d.String()
	}
}

func dbValues(config *Config) map[string]string {
	p := map[string]string{}
	setIfNotEmpty(p, "dbname", config.Name)
	setIfNotEmpty(p, "user", config.User)
	setIfNotEmpty(p, "host", config.Host)
	setIfNotEmpty(p, "port", config.Port)
	setIfNotEmpty(p, "sslmode", config.SSLMode)
	setIfPositive(p, "connect_timeout", config.ConnectionTimeout)
	setIfNotEmpty(p, "password", config.Password)
	setIfNotEmpty(p, "sslcert", config.SSLCertPath)
	setIfNotEmpty(p, "sslkey", config.SSLKeyPath)
	setIfNotEmpty(p, "sslrootcert", config.SSLRootCertPath)
	setIfNotEmpty(p, "pool_min_conns", config.PoolMinConnections)
	setIfNotEmpty(p, "pool_max_conns", config.PoolMaxConnections)
	setIfPositiveDuration(p, "pool_max_conn_lifetime", config.PoolMaxConnLife)
	setIfPositiveDuration(p, "pool_max_conn_idle_time", config.PoolMaxConnIdle)
	setIfPositiveDuration(p, "pool_health_check_period", config.PoolHealthCheck)
	return p
}
