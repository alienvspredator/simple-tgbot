package model

type Message struct {
	// Telegram's message ID
	TgMessageID int64
	UserID      int64
	ChatID      string
	Text        string
}
