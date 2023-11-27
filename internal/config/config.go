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
		ChatGPT  ChatGPT  `json:"chat_gpt"`
	}

	Postgres struct {
		URL string `json:"url"`
	}

	Telegram struct {
		Token string `json:"token"`
	}

	ChatGPT struct {
		Token string `json:"token"`
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
		ChatGPT: ChatGPT{
			Token: os.Getenv("TOKEN_CHAT_GPT"),
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
