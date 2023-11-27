package main

import (
	"github.com/goonec/business-tg-bot/internal/bot"
	"github.com/goonec/business-tg-bot/internal/config"
	"github.com/goonec/business-tg-bot/pkg/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed load config: %v", err)
	}

	if err := bot.Run(log, cfg); err != nil {
		log.Fatal("failed to run bot: %v", err)
	}
}
