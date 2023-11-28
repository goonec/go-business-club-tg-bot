package view

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/tgbot"
)

type viewResident struct {
	residentUsecase usecase.Resident
	log             *logger.Logger
}

func NewViewResident(residentUsecase usecase.Resident, log *logger.Logger) *viewResident {
	return &viewResident{
		residentUsecase: residentUsecase,
		log:             log,
	}
}

func (v *viewResident) ViewStart() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Hello")

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func (v *viewResident) ViewShowAllResident() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		fioMarkup, err := v.residentUsecase.GetAllFIOResident(context.Background())
		if err != nil {
			v.log.Info("residentUsecase.GetAllFIOResident: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, "<b>Список резидентов</b>")
		msg.ParseMode = tgbotapi.ModeHTML

		//msg := tgbotapi.EditMessageTextConfig{
		//	ParseMode: tgbotapi.ModeHTML,
		//	BaseEdit: tgbotapi.BaseEdit{
		//		ChatID:    update.Message.Chat.ID,
		//		MessageID: update.Message.MessageID,
		//	},
		//}

		msg.ReplyMarkup = fioMarkup

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
