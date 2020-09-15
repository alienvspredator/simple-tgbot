package tgbot

// Config is the configuration for the Bot components
type Config struct {
	TelegramToken string `env:"TG_TOKEN"`
	Debug         bool   `env:"LOG_DEBUG, default=false"`
}
