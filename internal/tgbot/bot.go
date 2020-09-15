// Package tgbot handles messages from Telegram users to the bot
package tgbot

import (
	"context"
	"fmt"
	"time"

	"github.com/alienvspredator/simple-tgbot/internal/logging"
	"github.com/alienvspredator/simple-tgbot/internal/serverenv"
	"github.com/yanzay/tbot/v2"
	"go.uber.org/zap"
)

type Bot struct {
	api *tbot.Server
}

// New builds a new bot application.
func New(env *serverenv.ServerEnv, config *Config) *Bot {
	api := tbot.New(config.TelegramToken /*, tbot.WithWebhook("", ":"+config.WebhookPort)*/)
	return &Bot{api: api}
}

// Serve runs the application loop. It stops on context done
func (b *Bot) Serve(ctx context.Context) error {
	log := logging.FromContext(ctx)

	go func() {
		<-ctx.Done()
		log.Info("stopping the receiving updates")
		b.api.Stop()
		log.Debug("the update receiver stopped")
	}()

	b.api.HandleMessage(".*", newEcho(ctx, b.api))

	if err := b.api.Start(); err != nil {
		return fmt.Errorf("start listening updates: %w", err)
	}

	return nil
}

func newEcho(ctx context.Context, api *tbot.Server) func(m *tbot.Message) {
	return func(m *tbot.Message) {
		log := logging.FromContext(ctx).With(
			"chat_id", m.Chat.ID,
			"incoming:message_id", m.MessageID,
		)

		log.Infow("got message", "text", m.Text)

		if err := api.Client().SendChatAction(m.Chat.ID, tbot.ActionTyping); err != nil {
			log.Errorw("send typing action", zap.Error(err))
		} else {
			time.Sleep(time.Second * 1)
		}

		if answer, err := api.Client().SendMessage(
			m.Chat.ID, m.Text, tbot.OptReplyToMessageID(m.MessageID),
		); err != nil {
			log.Errorw("send answer", zap.Error(err))
		} else {
			log.Infow("send answer", "answer:message_id", answer.MessageID, "text", answer.Text)
		}
	}
}
