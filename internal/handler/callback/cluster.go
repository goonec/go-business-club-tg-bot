package callback

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/localstore"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/tgbot"
)

type callbackBusinessCluster struct {
	businessCluster usecase.BusinessCluster
	residentUsecase usecase.Resident
	log             *logger.Logger

	store *localstore.Store
}

func NewCallbackBusinessCluster(businessCluster usecase.BusinessCluster, residentUsecase usecase.Resident, log *logger.Logger, store *localstore.Store) *callbackBusinessCluster {
	return &callbackBusinessCluster{
		businessCluster: businessCluster,
		residentUsecase: residentUsecase,
		log:             log,
		store:           store,
	}
}

func (c *callbackBusinessCluster) CallbackShowAllBusinessCluster() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		bcMarkup, err := c.businessCluster.GetAllBusinessCluster(ctx, "cluster_")
		if err != nil {
			c.log.Error("businessCluster.GetAllBusinessCluster: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Список кластеров")

		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = bcMarkup
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		return nil
	}
}

func (c *callbackBusinessCluster) CallbackGetIDCluster() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.CallbackQuery.Message.Chat.ID
		id := entity.FindID(update.CallbackData())
		if id == 0 {
			c.log.Error("entity.FindID: %v")
			return errors.New("zero id")
		}

		data, exist := c.store.Read(userID)
		if !exist {
			data = append(data, id)
			c.store.Set(data, userID)

			fioMarkup, err := c.residentUsecase.GetAllFIOResident(ctx, "getresident")
			if err != nil {
				c.store.Delete(userID)

				c.log.Error("residentUsecase.GetAllFIOResident: %v", err)
				handler.HandleError(bot, update, boterror.ParseErrToText(err))
			}

			msg := tgbotapi.NewMessage(update.FromChat().ID, `Выберите резедента`)
			msg.ParseMode = tgbotapi.ModeHTML

			msg.ReplyMarkup = fioMarkup

			if _, err := bot.Send(msg); err != nil {
				c.store.Delete(userID)
				return err
			}

			return nil
		}

		return nil
	}
}

func (c *callbackBusinessCluster) CallbackCreateClusterResident() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.CallbackQuery.Message.Chat.ID
		id := entity.FindID(update.CallbackData())
		if id == 0 {
			c.store.Delete(userID)
			c.log.Error("entity.FindID: %v")
			return errors.New("zero id")
		}

		data, exist := c.store.Read(userID)
		if !exist {
			c.store.Delete(userID)
			return errors.New("internal error")
		}

		data = append(data, id)
		err := c.businessCluster.CreateClusterResident(ctx, data[0].(int), data[1].(int)) //TODO переделать
		if err != nil {
			c.store.Delete(userID)
			c.log.Error("businessCluster.CreateClusterResident: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, `Кластер назначен резеденту`)

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		c.store.Delete(userID)
		return nil
	}
}

func (c *callbackBusinessCluster) CallbackDeleteCluster() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.CallbackQuery.Message.Chat.ID
		id := entity.FindID(update.CallbackData())
		if id == 0 {
			c.store.Delete(userID)
			c.log.Error("entity.FindID: %v")
			return errors.New("zero id")
		}

		err := c.businessCluster.DeleteCluster(ctx, id)
		if err != nil {
			c.log.Error("businessCluster.DeleteCluster: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, `Кластер удален успешно`)

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
