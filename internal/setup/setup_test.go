package setup_test

import (
	"context"
	"testing"

	"github.com/alienvspredator/simple-tgbot/internal/setup"
	"github.com/alienvspredator/simple-tgbot/internal/tgbot"
	"github.com/sethvargo/go-envconfig"
)

func TestSetupWith(t *testing.T) {
	t.Parallel()

	m := map[string]string{
		"TG_TOKEN": "TESTING_TOKEN",
	}
	lookuper := envconfig.MapLookuper(m)

	t.Run(
		"default", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			var config tgbot.Config
			env, err := setup.SetupWith(ctx, &config, lookuper)
			if err != nil {
				t.Fatal(err)
			}
			defer env.Close(ctx)

			if config.TelegramToken != "TESTING_TOKEN" {
				t.Errorf(
					"config.TelegramToken does not match, got = %s, want = %s",
					config.TelegramToken, m["TG_TOKEN"],
				)
			}
		},
	)
}
