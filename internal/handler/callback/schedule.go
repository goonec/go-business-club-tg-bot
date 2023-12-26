package callback

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/handler/tgbot"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
)

type callbackSchedule struct {
	scheduleUsecase usecase.Schedule
	log             *logger.Logger
}

func NewCallbackSchedule(scheduleUsecase usecase.Schedule, log *logger.Logger) *callbackSchedule {
	return &callbackSchedule{
		scheduleUsecase: scheduleUsecase,
		log:             log,
	}
}

func (c *callbackSchedule) CallbackGetSchedule() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		schedule, err := c.scheduleUsecase.GetSchedule(ctx)
		if err != nil {
			c.log.Info("scheduleUsecase.GetSchedule: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		schedulePhoto := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(schedule.PhotoFileID))

		msg := tgbotapi.NewPhoto(update.FromChat().ID, schedulePhoto.Media)
		//msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(handler.MainMenuButton))
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("%v")
			return err
		}
		return nil
	}
}
