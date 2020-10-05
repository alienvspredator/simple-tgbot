// This package is the simple Telegram Bot server
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/alienvspredator/simple-tgbot/internal/logging"
	"github.com/alienvspredator/simple-tgbot/internal/server"
	"github.com/alienvspredator/simple-tgbot/internal/setup"
	"github.com/alienvspredator/simple-tgbot/internal/tgbot"
	"github.com/sethvargo/go-signalcontext"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func main() {
	ctx, done := signalcontext.OnInterrupt()

	debug, _ := strconv.ParseBool(os.Getenv("LOG_DEBUG"))

	ctx = logging.WithLogger(ctx, logging.NewLogger(logging.WithDebugEnabled(debug)))

	err := realMain(ctx)
	done()

	log := logging.FromContext(ctx)

	if syncErr := log.Sync(); syncErr != nil {
		err = multierr.Append(err, syncErr)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func realMain(ctx context.Context) error {
	log := logging.FromContext(ctx)

	var config tgbot.Config
	ctx, env, err := setup.Setup(ctx, &config)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}
	defer func() {
		if err := env.Close(ctx); err != nil {
			log.Errorw("failed to close env", zap.Error(err))
		}
	}()

	bot := tgbot.New(env, &config)

	log.Info("starting bot")

	if err := server.ServeMetricsIfPrometheus(ctx); err != nil {
		return fmt.Errorf("serving metrics: %w", err)
	}
	return bot.Serve(ctx)
}
