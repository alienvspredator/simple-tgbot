// Package tgbot handles messages from Telegram users to the bot
package tgbot

import (
	"context"
	"fmt"
	"time"

	"github.com/alienvspredator/simple-tgbot/internal/logging"
	"github.com/alienvspredator/simple-tgbot/internal/metrics/metricsware"
	"github.com/alienvspredator/simple-tgbot/internal/serverenv"
	"github.com/alienvspredator/simple-tgbot/internal/tgbot/database"
	"github.com/alienvspredator/simple-tgbot/internal/tgbot/model"
	"github.com/yanzay/tbot/v2"
	"go.uber.org/zap"
)

type Bot struct {
	env    *serverenv.ServerEnv
	api    *tbot.Server
	db     *database.TgBotDB
	config *Config
}

// New builds a new bot application.
func New(env *serverenv.ServerEnv, config *Config) *Bot {
	return &Bot{
		env:    env,
		api:    tbot.New(config.TelegramToken),
		db:     database.New(env.Database()),
		config: config,
	}
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

	b.attachHandlers(ctx)

	if err := b.api.Start(); err != nil {
		return fmt.Errorf("start listening updates: %w", err)
	}

	return nil
}

func (b *Bot) attachHandlers(ctx context.Context) {
	b.api.HandleMessage("^/.*", b.EchoError(ctx))
	b.api.HandleMessage(".*", b.Echo(ctx))
}

func (b *Bot) EchoError(ctx context.Context) func(m *tbot.Message) {
	return func(m *tbot.Message) {
		log := logging.FromContext(ctx).With(
			"chat_id", m.Chat.ID,
			"incoming:message_id", m.MessageID,
		)
		metricsMW := metricsware.NewMiddleware()
		metricsMW.RecordIncomingMessage(ctx)

		log.Infow("got message", "text", m.Text)
		metricsMW.RecordFailedReply(ctx)
		metricsMW.RecordOutgoingMessage(ctx)
		if answer, err := b.api.Client().SendMessage(
			m.Chat.ID, "Я тебя не понимаю!", tbot.OptReplyToMessageID(m.MessageID),
		); err != nil {
			metricsMW.RecordSendFailure(ctx)
			log.Errorw("send answer", zap.Error(err))
		} else {
			log.Infow("send answer", "answer:message_id", answer.MessageID, "text", answer.Text)
		}
	}
}

func (b *Bot) Echo(ctx context.Context) func(m *tbot.Message) {
	return func(m *tbot.Message) {
		log := logging.FromContext(ctx).With(
			"chat_id", m.Chat.ID,
			"incoming:message_id", m.MessageID,
		)
		metricsMW := metricsware.NewMiddleware()
		metricsMW.RecordIncomingMessage(ctx)
		log.Infow("got message", "text", m.Text)

		if err := b.api.Client().SendChatAction(m.Chat.ID, tbot.ActionTyping); err != nil {
			log.Errorw("send typing action", zap.Error(err))
		} else {
			// Insert in separate goroutine
			go func() {
				if err := b.db.AddUserMessage(
					ctx, &model.Message{
						TgMessageID: int64(m.MessageID),
						UserID:      int64(m.From.ID),
						ChatID:      m.Chat.ID,
						Text:        m.Text,
					},
				); err != nil {
					metricsMW.RecordMessageSaveFailure(ctx)
					log.Errorw(
						"failed to save user message", "tg_message_id",
						m.MessageID, "user_id",
						m.From.ID, "chat_id",
						m.Chat.ID,
						"text", m.Text,
						zap.Error(err),
					)
					return
				}
			}()
			time.Sleep(time.Second * 1)
		}

		metricsMW.RecordOutgoingMessage(ctx)
		if answer, err := b.api.Client().SendMessage(
			m.Chat.ID, m.Text, tbot.OptReplyToMessageID(m.MessageID),
		); err != nil {
			metricsMW.RecordSendFailure(ctx)
			log.Errorw("send answer", zap.Error(err))
		} else {
			log.Infow("send answer", "answer:message_id", answer.MessageID, "text", answer.Text)
		}
	}
}
