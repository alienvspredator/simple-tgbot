package tgbot_test

import (
	"context"
	"testing"

	"github.com/alienvspredator/simple-tgbot/internal/serverenv"
	"github.com/alienvspredator/simple-tgbot/internal/tgbot"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run(
		"default", func(t *testing.T) {
			t.Parallel()

			env := serverenv.New(context.Background())
			config := &tgbot.Config{
				TelegramToken: "TESTING_TOKEN",
				Debug:         false,
			}

			bot := tgbot.New(env, config)
			if bot == nil {
				t.Fatal("bot was not created")
			}
		},
	)
}
