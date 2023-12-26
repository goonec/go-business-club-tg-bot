package view

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/handler/tgbot"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"time"
)

type viewSchedule struct {
	scheduleUsecase usecase.Schedule
	log             *logger.Logger

	transportChSchedule chan map[int64]map[string][]string
}

func NewViewSchedule(scheduleUsecase usecase.Schedule, log *logger.Logger, transportChSchedule chan map[int64]map[string][]string) *viewSchedule {
	return &viewSchedule{
		scheduleUsecase,
		log,
		transportChSchedule,
	}
}

func (v *viewSchedule) ViewCreateSchedule() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "[1] Отправьте фотографию расписания.")

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		go func() {
			subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			select {
			case d, ok := <-v.transportChSchedule:
				if ok {
					data := d[update.Message.From.ID]["/create_schedule"]
					if data == nil || len(data) == 0 {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, boterror.ParseErrToText(boterror.ErrInternalError))
						v.log.Error("ViewCreateSchedule: data == nil || len(data) == 0: %v", boterror.ErrInternalError)
						if _, err := bot.Send(msg); err != nil {
							v.log.Error("%v")
						}
						return
					}
					err := v.scheduleUsecase.CreateSchedule(context.Background(), data[0])
					if err != nil {
						v.log.Error("scheduleUsecase.CreateSchedule: %v", err)
						handler.HandleError(bot, update, boterror.ParseErrToText(err))
						return
					}

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Раписание добавлено успешно.")

					if _, err := bot.Send(msg); err != nil {
						v.log.Error("%v")
						handler.HandleError(bot, update, boterror.ParseErrToText(err))
					}
				}
				return
			case <-subCtx.Done():
				return
			}
		}()

		return nil
	}
}
