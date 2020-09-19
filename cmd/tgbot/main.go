// This package is the simple Telegram Bot server
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/alienvspredator/simple-tgbot/internal/logging"
	"github.com/alienvspredator/simple-tgbot/internal/setup"
	"github.com/alienvspredator/simple-tgbot/internal/tgbot"
	"github.com/sethvargo/go-signalcontext"
)

func main() {
	ctx, done := signalcontext.OnInterrupt()

	debug, _ := strconv.ParseBool(os.Getenv("LOG_DEBUG"))
	log := logging.NewLogger(debug)

	ctx = logging.WithLogger(ctx, log)

	err := realMain(ctx)
	done()

	if err != nil {
		log.Fatal(err)
	}
}

func realMain(ctx context.Context) error {
	log := logging.FromContext(ctx)

	var config tgbot.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}
	defer env.Close(ctx)

	bot := tgbot.New(env, &config)

	log.Info("starting bot")

	return bot.Serve(ctx)
}
