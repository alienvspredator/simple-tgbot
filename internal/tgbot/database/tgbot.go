package database

import (
	"context"
	"fmt"

	"github.com/alienvspredator/simple-tgbot/internal/database"
	"github.com/alienvspredator/simple-tgbot/internal/tgbot/model"
	"github.com/jackc/pgx/v4"
)

type TgBotDB struct {
	db *database.DB
}

func New(db *database.DB) *TgBotDB {
	return &TgBotDB{db: db}
}

func (db *TgBotDB) AddUserMessage(ctx context.Context, msg *model.Message) error {
	return db.db.InTx(
		ctx, pgx.Serializable, func(tx pgx.Tx) error {
			const q = `
			INSERT INTO
				received_messages
				(telegram_message_id, user_id, chat_id, message_text)
			VALUES
				($1, $2, $3, $4)
		`
			if _, err := tx.Exec(
				ctx, q, msg.TgMessageID, msg.UserID, msg.ChatID, msg.Text,
			); err != nil {
				return fmt.Errorf("saving message: %w", err)
			}

			return nil
		},
	)
}
