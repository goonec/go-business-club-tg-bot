package callback

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/tgbot"
)

type callbackBusinessCluster struct {
	businessCluster usecase.BusinessCluster
	log             *logger.Logger
}

func NewCallbackBusinessCluster(businessCluster usecase.BusinessCluster, log *logger.Logger) *callbackBusinessCluster {
	return &callbackBusinessCluster{
		businessCluster: businessCluster,
		log:             log,
	}
}

func (c *callbackBusinessCluster) CallbackShowAllBusinessCluster() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		bcMarkup, err := c.businessCluster.GetAllBusinessCluster(ctx)
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
