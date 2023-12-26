package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/handler/tgbot"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"time"
)

type callbackFeedback struct {
	feedbackUsecase     usecase.Feedback
	log                 *logger.Logger
	transportChFeedback chan map[int64][]string
}

func NewCallbackFeedback(feedbackUsecase usecase.Feedback, transportChFeedback chan map[int64][]string, log *logger.Logger) *callbackFeedback {
	return &callbackFeedback{
		feedbackUsecase:     feedbackUsecase,
		transportChFeedback: transportChFeedback,
		log:                 log,
	}
}

func (c *callbackFeedback) CallbackCreateFeedback(chatID int64) tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Пожалуйста отправьте сообщение для обратной связи. В сообщении укажите в виде данных номер и телеграмм, по которым можно будет связаться.")
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		go func(chatID int64) {
			subCtx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
			defer cancel()

			select {
			case d, exist := <-c.transportChFeedback:
				c.log.Info("", d, exist)
				data := d[update.CallbackQuery.Message.Chat.ID]
				if exist {
					if data == nil || len(data) == 0 {
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, boterror.ParseErrToText(boterror.ErrInternalError))
						c.log.Error("ViewCreateResidentPhoto: data == nil || len(data) == 0: %v", boterror.ErrInternalError)
						if _, err := bot.Send(msg); err != nil {
							c.log.Error("%v")
						}
						return
					}

					fb := &entity.Feedback{
						Message:    data[0],
						UsernameTG: update.CallbackQuery.From.UserName,
						Type:       "услуги",
					}

					_, err := c.feedbackUsecase.CreateFeedback(context.TODO(), fb)
					if err != nil {
						c.log.Error("feedbackUsecase.CreateFeedback: %v", err)
						handler.HandleError(bot, update, boterror.ParseErrToText(err))
						return
					}

					msg := tgbotapi.NewMessage(update.FromChat().ID, "Заявка оставлена успешно")
					if _, err := bot.Send(msg); err != nil {
						c.log.Error("failed to send message: %v", err)
						return
					}
					//c.sendFeedbackToAdmin(context.TODO(), feedback, chatID, bot)
					return
				}
			case <-subCtx.Done():
				return
			}
		}(chatID)
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
