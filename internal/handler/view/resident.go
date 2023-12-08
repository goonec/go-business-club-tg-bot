package view

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/tgbot"
	"strings"
	"sync"
	"time"
)

type viewResident struct {
	residentUsecase     usecase.Resident
	userUsecase         usecase.User
	log                 *logger.Logger
	transportCh         chan map[int64]map[string][]string
	transportСhResident chan map[int64]map[string][]string
}

func NewViewResident(residentUsecase usecase.Resident,
	userUsecase usecase.User,
	log *logger.Logger,
	transportCh chan map[int64]map[string][]string,
	transportСhResident chan map[int64]map[string][]string) *viewResident {
	return &viewResident{
		residentUsecase:     residentUsecase,
		userUsecase:         userUsecase,
		log:                 log,
		transportCh:         transportCh,
		transportСhResident: transportСhResident,
	}
}

func (v *viewResident) ViewAdminCommand() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		text := fmt.Sprintf("Доступные команды для администратора:\n/create_resident - создание резедента с его " +
			"фотографией и резюме\n/create_resident_photo - создание резедента только с фотографей\n" +
			"/notify - создание рассылки всем участникам бота\n/cancel - используется в случае отмены администраторской команды\n" +
			"/delete_resident - удаление резидента")
		msg := tgbotapi.NewMessage(update.FromChat().ID, text)

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func (v *viewResident) ViewShowAllResident() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		fioMarkup, err := v.residentUsecase.GetAllFIOResident(ctx, "")
		if err != nil {
			v.log.Error("residentUsecase.GetAllFIOResident: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, `<strong>Список резидентов</strong> 💼`)
		msg.ParseMode = tgbotapi.ModeHTML

		msg.ReplyMarkup = fioMarkup

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func (v *viewResident) ViewCreateResident() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "[1] Напишите ФИО, телеграм резидента,и его кластер."+
			" Должно получиться 5 слов, между которыми есть пробелы.")

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		go func() {
			subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			select {
			case d, ok := <-v.transportСhResident:
				data := d[update.Message.From.ID]["/create_resident"]
				if ok {
					if data == nil || len(data) == 0 {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, boterror.ParseErrToText(boterror.ErrInternalError))
						v.log.Error("ViewCreateResident: data == nil || len(data) == 0: %v", boterror.ErrInternalError)
						if _, err := bot.Send(msg); err != nil {
							v.log.Error("%v")
						}
						return
					}

					fioTg := strings.Split(data[0], " ")
					if len(fioTg) != 5 {
						handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrIncorrectAdminFirstInput))
						return
					}

					errStr := entity.IsFIOValid(fioTg[0], fioTg[1], fioTg[2])
					if len(errStr) != 0 {
						v.log.Error("entity.IsFIOValid: %v", errStr)
						handler.HandleError(bot, update, errStr)
						return
					}

					resident := &entity.Resident{
						BusinessCluster: entity.BusinessCluster{
							Name: fioTg[4],
						},
						FIO: entity.FIO{
							Firstname:  fioTg[0],
							Lastname:   fioTg[1],
							Patronymic: fioTg[2],
						},
						UsernameTG:   fioTg[3],
						ResidentData: data[1],
						PhotoFileID:  data[2],
					}

					err := v.residentUsecase.CreateResident(context.Background(), resident)
					if err != nil {
						v.log.Error("residentUsecase.CreateResident: %v", err)
						handler.HandleError(bot, update, boterror.ParseErrToText(err))
						return
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Резидент добавлен успешно.")

					if _, err := bot.Send(msg); err != nil {
						//return err
						v.log.Error("%v")
					}
				}
			case <-subCtx.Done():
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Истек срок создания резидента.")
				//if _, err := bot.Send(msg); err != nil {
				//	v.log.Error("%v", err)
				//}
				return
			}
		}()

		return nil
	}
}

func (v *viewResident) ViewCreateResidentPhoto() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "[1] Напишите ФИО, телеграм резидента,и его кластер."+
			" Должно получиться 5 слов, между которыми есть пробелы.")

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		go func() {
			subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			select {
			case d, ok := <-v.transportCh:
				data := d[update.Message.From.ID]["/create_resident_photo"]
				if ok {
					if data == nil || len(data) == 0 {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, boterror.ParseErrToText(boterror.ErrInternalError))
						v.log.Error("ViewCreateResidentPhoto: data == nil || len(data) == 0: %v", boterror.ErrInternalError)
						if _, err := bot.Send(msg); err != nil {
							v.log.Error("%v")
						}
						return
					}

					fioTg := strings.Split(data[0], " ")
					if len(fioTg) != 5 {
						handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrIncorrectAdminFirstInput))
						return
					}

					errStr := entity.IsFIOValid(fioTg[0], fioTg[1], fioTg[2])
					if len(errStr) != 0 {
						v.log.Error("entity.IsFIOValid: %v", errStr)
						handler.HandleError(bot, update, errStr)
						return
					}
					fmt.Println(data)

					resident := &entity.Resident{
						BusinessCluster: entity.BusinessCluster{
							Name: fioTg[4],
						},
						FIO: entity.FIO{
							Firstname:  fioTg[0],
							Lastname:   fioTg[1],
							Patronymic: fioTg[2],
						},
						UsernameTG:  fioTg[3],
						PhotoFileID: data[1],
					}

					err := v.residentUsecase.CreateResident(context.Background(), resident)
					if err != nil {
						v.log.Error("residentUsecase.CreateResident: %v", err)
						handler.HandleError(bot, update, boterror.ParseErrToText(err))
						return
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Резидент добавлен успешно.")

					if _, err := bot.Send(msg); err != nil {
						//return err
						v.log.Error("%v")
					}
				}
			case <-subCtx.Done():
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Истек срок создания резидента.")
				//if _, err := bot.Send(msg); err != nil {
				//	v.log.Error("%v", err)
				//}
				return
			}
		}()

		return nil
	}
}

func (v *viewResident) ViewCreateNotify() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "[1] Укажите сообщение для рассылки.")

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		go func() {
			var once sync.Once

			subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			select {
			case d, ok := <-v.transportCh:
				data := d[update.Message.From.ID]["/notify"]
				if ok {
					if data == nil || len(data) == 0 {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, boterror.ParseErrToText(boterror.ErrInternalError))
						v.log.Error("ViewCreateNotify: data == nil || len(data) == 0: %v", boterror.ErrInternalError)
						if _, err := bot.Send(msg); err != nil {
							v.log.Error("%v")
						}
						return
					}

					allID, err := v.userUsecase.GetAllUserID(context.Background())
					if err != nil {
						v.log.Error("userUsecase.GetAllUserID: %v", err)
						handler.HandleError(bot, update, boterror.ParseErrToText(err))
						return
					}

					for _, id := range allID {
						residentPhoto := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(data[1]))
						msg := tgbotapi.NewPhoto(id, residentPhoto.Media)
						msg.Caption = data[0]

						//msgText := tgbotapi.NewMessage(id, data[0])

						//if _, err := bot.Send(msgText); err != nil {
						//	v.log.Error("%v", err)
						//}

						if _, err := bot.Send(msg); err != nil {
							once.Do(func() {
								v.log.Error("%v :len(%d)", err, len([]rune(data[0])))
								errLongCap := errors.New(fmt.Sprintf("Количество символов в вашем сообщении: %d", len([]rune(data[0]))))
								err = errors.Join(err, errLongCap)
								handler.HandleError(bot, update, boterror.ParseErrToText(err))
							})
						}
					}
				}
			case <-subCtx.Done():
				return
			}
		}()

		return nil
	}
}

func (v *viewResident) ViewDeleteResident() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		fioMarkup, err := v.residentUsecase.GetAllFIOResident(ctx, "delete")
		if err != nil {
			v.log.Error("residentUsecase.GetAllFIOResident: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, "<b>Выберите резидента, которого нужно удалить:</b>")
		msg.ParseMode = tgbotapi.ModeHTML

		msg.ReplyMarkup = fioMarkup

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func (v *viewResident) ViewStartButton() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "<b>Список команд доступных для использования бота</b> ⏩")
		msg.ReplyMarkup = &handler.StartMenu
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
