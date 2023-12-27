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
	"time"
)

type viewService struct {
	serviceUsecase usecase.Service
	pptxUsecase    usecase.Pptx
	store          *localstore.Store

	transportPptx  chan map[int64]map[string][]string
	transportPhoto chan map[int64]map[string][]string
	log            *logger.Logger
}

func NewViewService(serviceUsecase usecase.Service,
	pptxUsecase usecase.Pptx,
	store *localstore.Store, log *logger.Logger,
	transportPptx chan map[int64]map[string][]string,
	transportPhoto chan map[int64]map[string][]string) *viewService {
	return &viewService{
		serviceUsecase: serviceUsecase,
		pptxUsecase:    pptxUsecase,
		store:          store,
		log:            log,
		transportPptx:  transportPptx,
		transportPhoto: transportPhoto,
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
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, `[1] Загрузите изображение для услуги`)
		if _, err := bot.Send(msg); err != nil {
			v.log.Error("%v", err)
			return err
		}

		go func() {
			subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			userID := update.Message.Chat.ID

			select {
			case d, ok := <-v.transportPhoto:
				dataPhoto := d[update.Message.From.ID]["/create_under_service"]
				if ok {
					if dataPhoto == nil || len(dataPhoto) == 0 {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, boterror.ParseErrToText(boterror.ErrInternalError))
						v.log.Error("ViewCreateResidentPhoto: data == nil || len(data) == 0: %v", boterror.ErrInternalError)
						if _, err := bot.Send(msg); err != nil {
							v.log.Error("%v")
						}
						return
					}

					type addServiceDescribe struct {
						Name     string `json:"under_service_name"`
						Describe string `json:"describe"`
					}

					data, exist := v.store.Read(userID)
					if !exist {
						args, err := parser.ParseJSON[addServiceDescribe](dataPhoto[1])
						if err != nil {
							v.log.Error("ParseJSON: %v", err)
							v.store.Delete(userID)
							handler.HandleError(bot, update, boterror.ParseErrToText(err))
							return
						}

						data = []interface{}{args.Name, args.Describe, dataPhoto[0]}
						v.store.Set(data, userID)

						sdMarkup, err := v.serviceUsecase.GetAllService(context.Background(), "create")
						if err != nil {
							v.log.Error("serviceUsecase.GetAllServiceDescribe: %v", err)
							v.store.Delete(userID)
							handler.HandleError(bot, update, boterror.ParseErrToText(err))
							return
						}

						msg := tgbotapi.NewMessage(update.FromChat().ID, `Выбирите услугу, которой нужно добавить раздел с описанием`)
						msg.ReplyMarkup = sdMarkup
						if _, err := bot.Send(msg); err != nil {
							v.log.Error("%v", err)
							v.store.Delete(userID)
							handler.HandleError(bot, update, boterror.ParseErrToText(err))
							return
						}
					} else {
						v.store.Delete(userID)

						msg := tgbotapi.NewMessage(update.FromChat().ID, `Произошла ошибка из-за прошлых операций. Попробуйте еще раз.`)
						if _, err := bot.Send(msg); err != nil {
							v.log.Error("%v", err)
							handler.HandleError(bot, update, boterror.ParseErrToText(err))
							return
						}
					}

				}
			case <-subCtx.Done():
				v.store.Delete(userID)
				return
			}

		}()
		return nil
	}
}

func (v *viewService) ViewCreatePptx() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Отправьте презентацию.")

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		go func() {
			subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			select {
			case d, ok := <-v.transportPptx:
				if ok {
					data := d[update.Message.From.ID]["/update_pptx"]

					err := v.pptxUsecase.UpdatePresentation(context.Background(), data[0])
					if err != nil {
						v.log.Error("%v", err)
						handler.HandleError(bot, update, boterror.ParseErrToText(err))
						return
					}

					msg := tgbotapi.NewMessage(update.FromChat().ID, `Презентация загружена успешно.`)
					if _, err := bot.Send(msg); err != nil {
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

func (v *viewService) ViewDeleteService() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		sMarkup, err := v.serviceUsecase.GetAllService(ctx, "delete")
		if err != nil {
			v.log.Error("serviceUsecase.GetAllService: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, `Выбирите сервис, который нужно удалить.`)
		msg.ReplyMarkup = sMarkup

		if _, err := bot.Send(msg); err != nil {
			v.log.Error("%v", err)
			return err
		}

		return nil
	}
}

func (v *viewService) ViewDeleteUnderService() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		sdMarkup, err := v.serviceUsecase.GetAllServiceDescribe(ctx, "descdelete")
		if err != nil {
			v.log.Error("serviceUsecase.GetAllService: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, `Выбирите раздел сервиса, который нужно удалить.`)
		msg.ReplyMarkup = sdMarkup

		if _, err := bot.Send(msg); err != nil {
			v.log.Error("%v", err)
			return err
		}

		return nil
	}
}
