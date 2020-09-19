package tgbot

import (
	"github.com/alienvspredator/simple-tgbot/internal/database"
	"github.com/alienvspredator/simple-tgbot/internal/secrets"
)

// Config is the configuration for the Bot components
type Config struct {
	Database      database.Config
	SecretManager secrets.Config

	TelegramToken string `env:"TG_TOKEN"`
	Debug         bool   `env:"LOG_DEBUG, default=false"`
	WebhookPort   string `env:"PORT"`
}

func (c *Config) SecretManagerConfig() *secrets.Config {
	return &c.SecretManager
}

func (c *Config) DatabaseConfig() *database.Config {
	return &c.Database
}
