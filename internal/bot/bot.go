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

	transportCh := make(chan map[int64]map[string][]string, 1)
	transportСhResident := make(chan map[int64]map[string][]string, 1)

	residentRepo := repo.NewResidentRepository(psql)
	userRepo := repo.NewUserRepository(psql)
	businessClusterRepo := repo.NewBusinessClusterRepository(psql)
	businessClusterResidentRepo := repo.NewBusinessClusterResidentRepository(psql)

	residentUsecase := usecase.NewResidentUsecase(residentRepo, businessClusterRepo, businessClusterResidentRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	businessClusterUsecase := usecase.NewBusinessClusterUsecase(businessClusterRepo)

	residentView := view.NewViewResident(residentUsecase, userUsecase, log, transportCh, transportСhResident)

	residentCallback := callback.NewCallbackResident(residentUsecase, log)
	businessClusterCallback := callback.NewCallbackBusinessCluster(businessClusterUsecase, log)

	newBot := tgbot.NewBot(bot, log, openaiRequest, userUsecase, transportCh, transportСhResident, cfg.Chat.ChatID)
	newBot.RegisterCommandView("admin", middleware.AdminMiddleware(cfg.Chat.ChatID, residentView.ViewAdminCommand()))
	newBot.RegisterCommandView("create_resident", middleware.AdminMiddleware(cfg.Chat.ChatID, residentView.ViewCreateResident()))
	newBot.RegisterCommandView("create_resident_photo", middleware.AdminMiddleware(cfg.Chat.ChatID, residentView.ViewCreateResidentPhoto()))
	newBot.RegisterCommandView("notify", middleware.AdminMiddleware(cfg.Chat.ChatID, residentView.ViewCreateNotify()))
	newBot.RegisterCommandView("delete_resident", middleware.AdminMiddleware(cfg.Chat.ChatID, residentView.ViewDeleteResident()))

	newBot.RegisterCommandView("start", residentView.ViewStartButton())
	newBot.RegisterCommandView("resident_list", residentView.ViewShowAllResident())

	//newBot.RegisterCommandCallback("stop_chat_gpt", residentCallback.CallbackStopChatGPT())
	newBot.RegisterCommandCallback("resident", residentCallback.CallbackShowAllResident())
	newBot.RegisterCommandCallback("chat_gpt", residentCallback.CallbackStartChatGPT())

	newBot.RegisterCommandCallback("main_menu", residentCallback.CallbackStartButton())

	newBot.RegisterCommandCallback("cluster", residentCallback.CallbackShowResidentByCluster())
	newBot.RegisterCommandCallback("allcluster", businessClusterCallback.CallbackShowAllBusinessCluster())
	newBot.RegisterCommandCallback("fiodelete", residentCallback.CallbackDeleteResident())
	newBot.RegisterCommandCallback("fio", residentCallback.CallbackGetResident())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := newBot.Run(ctx); err != nil {
		log.Error("failed to run tgbot: %v", err)
	}

	return nil
}
