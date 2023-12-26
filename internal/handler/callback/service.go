package callback

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/tgbot"
)

type callbackService struct {
	serviceUsecase usecase.Service
	log            *logger.Logger
}

func NewCallbackService(serviceUsecase usecase.Service, log *logger.Logger) *callbackService {
	return &callbackService{
		serviceUsecase: serviceUsecase,
		log:            log,
	}
}

func (c *callbackService) ViewShowAllService() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		serviceMarkup, err := c.serviceUsecase.GetAllService(ctx)
		if err != nil {
			c.log.Error("serviceUsecase.GetAllService: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, "Выбирете услугу")

		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = &serviceMarkup
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		return nil
	}
}

func (c *callbackService) ViewShowAllServiceDescribe() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		id := entity.FindID(update.CallbackData())
		if id == 0 {
			c.log.Error("entity.FindID: %v")
			handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrIncorrectCallbackData))
			return nil
		}

		serviceDescribeMarkup, err := c.serviceUsecase.GetAllServiceDescribe(ctx, id, "describe")
		if err != nil {
			c.log.Error("serviceUsecase.GetAllServiceDescribe: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Название услуг")
		msg.ReplyMarkup = serviceDescribeMarkup

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
