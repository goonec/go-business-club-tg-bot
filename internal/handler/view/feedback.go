package view

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/handler/tgbot"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
)

type viewFeedback struct {
	feedbackUsecase usecase.Feedback
	log             *logger.Logger
}

func NewViewFeedback(feedbackUsecase usecase.Feedback, log *logger.Logger) *viewFeedback {
	return &viewFeedback{
		feedbackUsecase: feedbackUsecase,
		log:             log,
	}
}

func (v *viewFeedback) CallbackGetFeedback() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		feedback, err := v.feedbackUsecase.GetAllFeedback(ctx)
		if err != nil {
			v.log.Error("feedbackUsecase.GetAllFeedback: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		feedbackByte, _ := json.MarshalIndent(feedback, "", "")
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s", string(feedbackByte)))
		if _, err := bot.Send(msg); err != nil {
			v.log.Error("failed to send message: %v", err)
			return err
		}
		return nil
	}
}
