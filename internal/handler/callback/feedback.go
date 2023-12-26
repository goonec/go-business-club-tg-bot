package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/handler/tgbot"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
)

type callbackFeedback struct {
	feedbackUsecase usecase.Feedback
	log             *logger.Logger
}

func NewCallbackFeedback(feedbackUsecase usecase.Feedback, log *logger.Logger) *callbackFeedback {
	return &callbackFeedback{
		feedbackUsecase: feedbackUsecase,
		log:             log,
	}
}

func (c *callbackFeedback) CallbackCreateFeedback() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

		return nil
	}
}

func (c *callbackFeedback) sendFeedbackToAdmin(ctx context.Context, feedback *entity.Feedback, channelID int64, bot *tgbotapi.BotAPI) {
	select {
	case <-ctx.Done():
		c.log.Error("cancel context")
		return
	default:
		admins, err := bot.GetChatAdministrators(
			tgbotapi.ChatAdministratorsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: channelID,
				},
			})

		if err != nil {
			c.log.Error("GetChatAdministrators: %v", err)
			return
		}

		for _, admin := range admins {
			text := fmt.Sprintf("Пришла обратная связь:\n %#v", feedback)
			msg := tgbotapi.NewMessage(admin.User.ID, text)
			if _, err := bot.Send(msg); err != nil {
				c.log.Error("%v", err)
				return
			}

		}
	}
}
