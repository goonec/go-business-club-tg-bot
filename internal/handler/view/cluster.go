package view

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/tgbot"
)

type viewCluster struct {
	clusterUsecase usecase.BusinessCluster
	log            *logger.Logger
}

func NewViewCluster(clusterUsecase usecase.BusinessCluster, log *logger.Logger) *viewCluster {
	return &viewCluster{
		clusterUsecase: clusterUsecase,
		log:            log,
	}
}

func ParseJSON[T any](src string) (T, error) {
	var args T

	if err := json.Unmarshal([]byte(src), &args); err != nil {
		return *(new(T)), err
	}

	return args, nil
}

func (v *viewCluster) ViewCreateCluster() tgbot.ViewFunc {
	type addClusterArgs struct {
		Cluster string `json:"cluster"`
	}
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		args, err := ParseJSON[addClusterArgs](update.Message.CommandArguments())
		if err != nil {
			v.log.Error("ParseJSON: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		err = v.clusterUsecase.Create(ctx, args.Cluster)
		if err != nil {
			v.log.Error("clusterRepo.Create: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, `Кластер добавлен успешно.`)
		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
