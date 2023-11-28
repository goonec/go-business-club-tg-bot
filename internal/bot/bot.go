package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/config"
	"github.com/goonec/business-tg-bot/internal/handler/middleware"
	"github.com/goonec/business-tg-bot/internal/handler/view"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/openai"
	"github.com/goonec/business-tg-bot/pkg/tgbot"
	"os"
	"os/signal"
	"syscall"
)

func Run(log *logger.Logger, cfg *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("failed to load token %v", err)
	}

	bot.Debug = false

	log.Info("Authorized on account %s", bot.Self.UserName)

	openaiRequest := openai.NewOpenAIConnect(cfg.OpenAI.Token)
	residentView := view.NewViewResident(nil, log)

	newBot := tgbot.NewBot(bot, log, openaiRequest)
	newBot.RegisterCommandView("start", middleware.AdminMiddleware(cfg.Chat.ChatID, residentView.ViewStart()))
	//newBot.RegisterCommandView("notify", middleware.AdminMiddleware(cfg.Chat.ChatID))
	//newBot.RegisterCommandView("resident_list")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := newBot.Run(ctx); err != nil {
		log.Error("failed to run tgbot: %v", err)
	}

	return nil
}
