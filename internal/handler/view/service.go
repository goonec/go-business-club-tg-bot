package view

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/handler/tgbot"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/localstore"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/parser"
	"github.com/goonec/business-tg-bot/pkg/postgres"
	"time"
)

type viewService struct {
	serviceUsecase usecase.Service
	store          *localstore.Store

	transportPptx chan map[int64]map[string][]string
	log           *logger.Logger
	pg            *postgres.Postgres
}

func NewViewService(serviceUsecase usecase.Service, store *localstore.Store, log *logger.Logger, transportPptx chan map[int64]map[string][]string, pg *postgres.Postgres) *viewService {
	return &viewService{
		serviceUsecase: serviceUsecase,
		store:          store,
		log:            log,
		transportPptx:  transportPptx,
		pg:             pg,
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
			v.log.Error("%v", err)
			return err
		}

		return nil
	}
}

func (v *viewService) ViewCreateUnderService() tgbot.ViewFunc {
	type addServiceDescribe struct {
		Name     string `json:"under_service_name"`
		Describe string `json:"describe"`
	}
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.Message.Chat.ID
		data, exist := v.store.Read(userID)
		if !exist {
			args, err := parser.ParseJSON[addServiceDescribe](update.Message.CommandArguments())
			if err != nil {
				v.log.Error("ParseJSON: %v", err)
				v.store.Delete(userID)
				handler.HandleError(bot, update, boterror.ParseErrToText(err))
				return nil
			}

			data = []interface{}{args.Name, args.Describe}
			v.store.Set(data, userID)

			sdMarkup, err := v.serviceUsecase.GetAllService(ctx, "create")
			if err != nil {
				v.log.Error("serviceUsecase.GetAllServiceDescribe: %v", err)
				v.store.Delete(userID)
				handler.HandleError(bot, update, boterror.ParseErrToText(err))
				return nil
			}

			msg := tgbotapi.NewMessage(update.FromChat().ID, `Выбирите услугу, которой нужно добавить раздел с описанием`)
			msg.ReplyMarkup = sdMarkup
			if _, err := bot.Send(msg); err != nil {
				v.log.Error("%v", err)
				v.store.Delete(userID)
				return err
			}
		} else {
			v.store.Delete(userID)

			msg := tgbotapi.NewMessage(update.FromChat().ID, `Произошла ошибка из-за прошлых операций. Попробуйте еще раз.`)
			if _, err := bot.Send(msg); err != nil {
				v.log.Error("%v", err)
				return err
			}
		}
		return nil
	}
}

func (v *viewService) ViewCreatePhotoServiceDescribe() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

		return nil
	}
}

func (v *viewService) ViewCreatePptx() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Отправьте презентацию")

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		go func() {
			subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			select {
			case d, ok := <-v.transportPptx:
				v.log.Info("", d, ok)
				if ok {
					data := d[update.Message.From.ID]["/create_pptx"]

					_, err := v.pg.Pool.Exec(context.Background(), "insert into pptx (pptx_file_id) values ($1)", data[0])
					if err != nil {
						v.log.Error("%v", err)
						return
					}
				}
				return
			case <-subCtx.Done():
				return
			}

			return
		}()
		return nil
	}
}
