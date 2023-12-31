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
	transportChFeedback chan map[int64][]interface{}
}

func NewCallbackFeedback(feedbackUsecase usecase.Feedback, transportChFeedback chan map[int64][]interface{}, log *logger.Logger) *callbackFeedback {
	return &callbackFeedback{
		feedbackUsecase:     feedbackUsecase,
		transportChFeedback: transportChFeedback,
		log:                 log,
	}
}

func (c *callbackFeedback) CallbackCreateFeedback(chatID int64) tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Пожалуйста отправьте сообщение для обратной связи. В сообщении укажите в виде данных номер и телеграмм, по которым можно будет связаться.\n\n"+
			"В случае, если вы передумали отправлять сообщение, воспользуйтесь командой - /exit.")
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		go func() {
			subCtx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
			defer cancel()

			select {
			case d, exist := <-c.transportChFeedback:
				c.log.Info("feedback", d, exist)
				if exist {
					for key, value := range d {
						u := value[1].(tgbotapi.Update)
						fb := &entity.Feedback{
							Message:    value[0].(string),
							UsernameTG: u.Message.From.UserName,
							Type:       value[2].(string),
						}

						_, err := c.feedbackUsecase.CreateFeedback(context.TODO(), fb)
						if err != nil {
							c.log.Error("feedbackUsecase.CreateFeedback: %v", err)
							handler.HandleError(bot, update, boterror.ParseErrToText(err))
							return
						}

						msg := tgbotapi.NewMessage(key, "Заявка оставлена успешно")
						if _, err := bot.Send(msg); err != nil {
							c.log.Error("failed to send message: %v", err)
							return
						}
						//c.sendFeedbackToAdmin(context.TODO(), feedback, chatID, bot)
						return
					}
				}
			case <-subCtx.Done():
				return
			}
		}()
		return nil
	}
}

func (c *callbackFeedback) CallbackMembershipRequest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Пожалуйста, пришлите нам ваше ФИО, номер телефона и вид деятельности. Мы свяжемся с вами в ближайшее время!\n\n"+
			"В случае, если вы передумали отправлять сообщение, воспользуйтесь командой - /exit.")
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		go func() {
			subCtx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
			defer cancel()

			select {
			case d, exist := <-c.transportChFeedback:
				c.log.Info("membership request", d, exist)
				if exist {
					for key, value := range d {
						u := value[1].(tgbotapi.Update)
						fb := &entity.Feedback{
							Message:    value[0].(string),
							UsernameTG: u.Message.From.UserName,
							Type:       value[2].(string),
						}

						_, err := c.feedbackUsecase.CreateFeedback(context.TODO(), fb)
						if err != nil {
							c.log.Error("feedbackUsecase.CreateFeedback: %v", err)
							handler.HandleError(bot, update, boterror.ParseErrToText(err))
							return
						}

						msg := tgbotapi.NewMessage(key, "Заявка оставлена успешно")
						if _, err := bot.Send(msg); err != nil {
							c.log.Error("failed to send message: %v", err)
							return
						}
						//c.sendFeedbackToAdmin(context.TODO(), feedback, chatID, bot)
						return
					}
				}
			case <-subCtx.Done():
				return
			}
		}()

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
