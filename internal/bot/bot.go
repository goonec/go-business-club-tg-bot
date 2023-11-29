package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/config"
	"github.com/goonec/business-tg-bot/internal/handler/callback"
	"github.com/goonec/business-tg-bot/internal/handler/middleware"
	"github.com/goonec/business-tg-bot/internal/handler/view"
	"github.com/goonec/business-tg-bot/internal/repo"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/openai"
	"github.com/goonec/business-tg-bot/pkg/postgres"
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

	psql, err := postgres.New(context.Background(), 5, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}
	defer psql.Close()

	log.Info("Authorized on account %s", bot.Self.UserName)

	openaiRequest := openai.NewOpenAIConnect(cfg.OpenAI.Token)

	transportCh := make(chan []string)

	residentRepo := repo.NewResidentRepository(psql)
	residentUsecase := usecase.NewResidentUsecase(residentRepo)
	residentView := view.NewViewResident(residentUsecase, log, transportCh)

	residentCallback := callback.NewCallbackResident(residentUsecase, log)

	newBot := tgbot.NewBot(bot, log, openaiRequest, transportCh)
	newBot.RegisterCommandView("start", middleware.AdminMiddleware(cfg.Chat.ChatID, residentView.ViewStart()))
	newBot.RegisterCommandView("create_resident", middleware.AdminMiddleware(cfg.Chat.ChatID, residentView.ViewCreateResident()))
	newBot.RegisterCommandView("resident_list", residentView.ViewShowAllResident())

	newBot.RegisterCommandCallback("fio", residentCallback.CallbackGetResident())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := newBot.Run(ctx); err != nil {
		log.Error("failed to run tgbot: %v", err)
	}

	return nil
}
