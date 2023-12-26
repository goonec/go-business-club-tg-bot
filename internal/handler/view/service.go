package view

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/parser"
	"github.com/goonec/business-tg-bot/pkg/tgbot"
)

type viewService struct {
	serviceUsecase usecase.Service
	log            *logger.Logger
}

func NewViewService(serviceUsecase usecase.Service, log *logger.Logger) *viewService {
	return &viewService{
		serviceUsecase: serviceUsecase,
		log:            log,
	}
}

func (v *viewService) ViewCreateService() tgbot.ViewFunc {
	type addServiceArgs struct {
		ServiceName string `json:"service_name"`
	}
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		args, err := parser.ParseJSON[addServiceArgs](update.Message.CommandArguments())
		if err != nil {
			v.log.Error("ParseJSON: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		err = v.serviceUsecase.CreateService(ctx, args.ServiceName)
		if err != nil {
			v.log.Error("serviceUsecase.CreateService: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, `Раздел услуги добавлен успешно.`)
		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
