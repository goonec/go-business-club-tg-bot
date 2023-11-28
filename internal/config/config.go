package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type (
	Config struct {
		Postgres Postgres `json:"postgres"`
		Telegram Telegram `json:"telegram"`
		OpenAI   OpenAI   `json:"open_ai"`
		Chat     Chat     `json:"chat"`
	}

	Postgres struct {
		URL string `json:"url"`
	}

	Telegram struct {
		Token string `json:"token"`
	}

	OpenAI struct {
		Token string `json:"token"`
	}

	Chat struct {
		ChatID int64 `json:"chat"`
	}
)

func New() (*Config, error) {
	err := godotenv.Load("configs/bot.env")
	if err != nil {
		return nil, err
	}

	config := &Config{
		Postgres: Postgres{
			URL: os.Getenv("POSTGRES_URL"),
		},
		Telegram: Telegram{
			Token: os.Getenv("TOKEN_TG"),
		},
		OpenAI: OpenAI{
			Token: os.Getenv("TOKEN_CHAT_GPT"),
		},
		Chat: Chat{
			ChatID: int64(parseEnvInt(os.Getenv("CHANNEL_ID"))),
		},
	}

	return config, nil
}

func parseEnvInt(value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return intValue
}
