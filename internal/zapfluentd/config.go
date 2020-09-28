package zapfluentd

type Config struct {
	Host string `env:"FLUENT_HOST"`
	Port int    `env:"FLUENT_PORT, default=24224"`
}
