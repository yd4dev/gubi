package config

import (
	"log/slog"
	"os"

	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

type Config struct {
	Token    string
	DevGuild snowflake.ID
}

func Load() *Config {
	if err := godotenv.Load(); err == nil {
		slog.Info(".env file found, reading variables from .env.")
	}

	return &Config{
		Token:    os.Getenv("BOT_TOKEN"),
		DevGuild: snowflake.GetEnv("DEV_GUILD"),
	}
}
