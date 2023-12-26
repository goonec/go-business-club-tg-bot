package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/localstore"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/tgbot"
)

type callbackService struct {
	serviceUsecase usecase.Service
	store          *localstore.Store
	log            *logger.Logger
}

func NewCallbackService(serviceUsecase usecase.Service, store *localstore.Store, log *logger.Logger) *callbackService {
	return &callbackService{
		serviceUsecase: serviceUsecase,
		store:          store,
		log:            log,
	}
}

func (c *callbackService) ViewShowAllService() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		serviceMarkup, err := c.serviceUsecase.GetAllService(ctx, "")
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

		serviceDescribeMarkup, err := c.serviceUsecase.GetAllServiceDescribeByServiceID(ctx, id, "describe")
		if err != nil {
			c.log.Error("serviceUsecase.GetAllServiceDescribeByServiceID: %v", err)
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

func (c *callbackService) ViewShowServiceInfo() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		id := entity.FindID(update.CallbackData())
		if id == 0 {
			c.log.Error("entity.FindID: %v")
			handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrIncorrectCallbackData))
			return nil
		}

		serviceDescribe, err := c.serviceUsecase.Get(ctx, id)
		if err != nil {
			c.log.Error("serviceUsecase.Get: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf(
			"Название услуги: %s"+
				"\n\n"+
				"Название раздела: %s\n"+
				"Описание: %s", serviceDescribe.Service.Name, serviceDescribe.Name, serviceDescribe.Describe))

		serviceDescribeMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(handler.FeedbackButton),
			tgbotapi.NewInlineKeyboardRow(handler.MainMenuButton))

		msg.ReplyMarkup = &serviceDescribeMarkup

		if _, err := bot.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
