package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/config"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/openai"
)

func Run(log *logger.Logger, cfg *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("failed to load token %v", err)
	}

	bot.Debug = false

	log.Info("Authorized on account %s", bot.Self.UserName)

	openaiRequest := openai.NewOpenAIConnect(cfg.OpenAI.Token, "")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			switch update.Message.Command() {
			case "start":
				msg.Text = "Бизнес клуб бот"

				_, err = bot.Send(msg)
				if err != nil {
					log.Error("failed to send message in /start %v", err)
				}
			default:
				openaiResponse, err := openaiRequest.ResponseGPT(update.Message.Text)
				if err != nil {
					log.Error("failed to get response from GPT: %v", err)
				}

				msg.Text = openaiResponse

				_, err = bot.Send(msg)
				if err != nil {
					log.Error("failed to send message from ChatGPT %v", err)
				}
			}
		}
	}

	return nil
}
