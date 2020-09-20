package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/alienvspredator/simple-tgbot/internal/database"
	"github.com/alienvspredator/simple-tgbot/internal/logging"
	"github.com/alienvspredator/simple-tgbot/internal/setup"
	"github.com/golang-migrate/migrate/v4"
	"github.com/markbates/pkger"
	"github.com/sethvargo/go-signalcontext"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/pkger"
	_ "github.com/lib/pq"
)

func main() {
	ctx, done := signalcontext.OnInterrupt()
	log := logging.NewLogger(false)

	ctx = logging.WithLogger(ctx, log)

	err := realMain(ctx)
	done()

	if err != nil {
		log.Fatal(err)
	}
}

func realMain(ctx context.Context) error {
	log := logging.FromContext(ctx)

	var config database.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		return fmt.Errorf("setup database: %w", err)
	}
	defer env.Close(ctx)

	pkger.Include("/migrations")
	m, err := migrate.New("pkger:///migrations", config.ConnectionURL())
	if err != nil {
		return fmt.Errorf("create migrate: %w", err)
	}
	m.Log = newLogger(ctx)

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("already up to date")
		} else {
			return fmt.Errorf("run migration up: %w", err)
		}
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return fmt.Errorf("migrate source: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("migrate database: %w", dbErr)
	}

	log.Debugw("finished migrations up")
	return nil
}

type logger struct {
	zap *zap.SugaredLogger
}

func newLogger(ctx context.Context) *logger {
	return &logger{
		zap: logging.FromContext(ctx).Named("migrate"),
	}
}

func (l *logger) Printf(arg string, vars ...interface{}) {
	l.zap.Infof(arg, vars...)
}

func (l *logger) Verbose() bool {
	return true
}
