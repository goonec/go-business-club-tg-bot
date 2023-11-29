package view

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/tgbot"
	"strings"
)

type viewResident struct {
	residentUsecase usecase.Resident
	log             *logger.Logger
	transportCh     chan []string
}

func NewViewResident(residentUsecase usecase.Resident, log *logger.Logger, transportCh chan []string) *viewResident {
	return &viewResident{
		residentUsecase: residentUsecase,
		log:             log,
		transportCh:     transportCh,
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
		fioMarkup, err := v.residentUsecase.GetAllFIOResident(ctx)
		if err != nil {
			v.log.Error("residentUsecase.GetAllFIOResident: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, "<b>Список резидентов</b>")
		msg.ParseMode = tgbotapi.ModeHTML

		msg.ReplyMarkup = fioMarkup

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

// tg, fio, describe, photo
func (v *viewResident) ViewCreateResident() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "[1] Напишите ФИО и телеграм резидента."+
			" Должно получиться 4 слова, между которыми есть пробелы.")

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		go func() {
			if data, ok := <-v.transportCh; ok {
				fioTg := strings.Split(data[0], " ")
				if len(fioTg) != 4 {
					handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrIncorrectAdminFirstInput))
					return
				}

				err := entity.IsFIOValid(fioTg[0], fioTg[1], fioTg[2])
				if err != nil {
					v.log.Error("entity.IsFIOValid: %v", err)
					handler.HandleError(bot, update, err.Error())
					return
				}

				resident := &entity.Resident{
					FIO: entity.FIO{
						Firstname:  fioTg[0],
						Lastname:   fioTg[1],
						Patronymic: fioTg[2],
					},
					UsernameTG:   fioTg[3],
					ResidentData: data[1],
					PhotoFileID:  data[2],
				}
				err = v.residentUsecase.CreateResident(context.Background(), resident)
				if err != nil {
					v.log.Error("residentUsecase.CreateResident: %v", err)
					handler.HandleError(bot, update, boterror.ParseErrToText(err))
					return
				}
				msg := tgbotapi.NewMessage(update.FromChat().ID, "Резидент добавлен успешно.")

				if _, err := bot.Send(msg); err != nil {
					//return err
					v.log.Error("%v")
				}

			}
		}()

		return nil
	}
}
