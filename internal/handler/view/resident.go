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
	"time"
)

type viewResident struct {
	residentUsecase usecase.Resident
	log             *logger.Logger
	transportCh     chan map[int64][]string
}

func NewViewResident(residentUsecase usecase.Resident, log *logger.Logger, transportCh chan map[int64][]string) *viewResident {
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

func (v *viewResident) ViewCreateResident() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "[1] Напишите ФИО и телеграм резидента."+
			" Должно получиться 4 слова, между которыми есть пробелы.")

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		go func() {
			subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			select {
			case d, ok := <-v.transportCh:
				data := d[update.Message.From.ID]
				if ok {
					fioTg := strings.Split(data[0], " ")
					if len(fioTg) != 4 {
						handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrIncorrectAdminFirstInput))
						return
					}

					errStr := entity.IsFIOValid(fioTg[0], fioTg[1], fioTg[2])
					if len(errStr) != 0 {
						v.log.Error("entity.IsFIOValid: %v", errStr)
						handler.HandleError(bot, update, errStr)
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
					err := v.residentUsecase.CreateResident(context.Background(), resident)
					if err != nil {
						v.log.Error("residentUsecase.CreateResident: %v", err)
						handler.HandleError(bot, update, boterror.ParseErrToText(err))
						return
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Резидент добавлен успешно.")

					if _, err := bot.Send(msg); err != nil {
						//return err
						v.log.Error("%v")
					}
				}
			case <-subCtx.Done():
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Истек срок создания резидента.")
				//if _, err := bot.Send(msg); err != nil {
				//	v.log.Error("%v", err)
				//}
				return
			}
		}()

		return nil
	}
}
