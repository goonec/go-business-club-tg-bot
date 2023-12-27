package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/handler"
	"github.com/goonec/business-tg-bot/internal/handler/tgbot"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/localstore"
	"github.com/goonec/business-tg-bot/pkg/logger"
)

type callbackResident struct {
	residentUsecase usecase.Resident
	log             *logger.Logger

	store *localstore.Store
}

func NewCallbackResident(residentUsecase usecase.Resident, log *logger.Logger, store *localstore.Store) *callbackResident {
	return &callbackResident{
		residentUsecase: residentUsecase,
		log:             log,
		store:           store,
	}
}

func (c *callbackResident) CallbackGetResident() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		id := entity.FindID(update.CallbackData())
		if id == 0 {
			c.log.Error("entity.FindID == 0")
			//return boterror.ErrIncorrectCallbackData
			handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrIncorrectCallbackData))
			return nil
		}

		resident, err := c.residentUsecase.GetResident(ctx, id)
		if err != nil {
			c.log.Info("residentUsecase.GetResident: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}
		userID := update.CallbackQuery.Message.Chat.ID

		residentPhoto := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(resident.PhotoFileID))

		msg := tgbotapi.NewPhoto(userID, residentPhoto.Media)
		text := fmt.Sprintf("%s %s\n\n%s", resident.FIO.Firstname, resident.FIO.Lastname, resident.ResidentData)

		msgText := tgbotapi.NewMessage(userID, text)

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		if _, err := bot.Send(msgText); err != nil {
			return err
		}

		return nil
	}
}

func (c *callbackResident) CallbackDeleteResident() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		id := entity.FindID(update.CallbackData())
		if id == 0 {
			c.log.Error("entity.FindID == 0")
			handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrIncorrectCallbackData))
		}

		err := c.residentUsecase.DeleteResident(ctx, id)
		if err != nil {
			c.log.Info("residentUsecase.DeleteResident: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
		}
		userID := update.CallbackQuery.Message.Chat.ID

		msg := tgbotapi.NewMessage(userID, "Резидент удален успешно.")

		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		return nil
	}
}

func (c *callbackResident) CallbackShowAllResident() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		fioMarkup, err := c.residentUsecase.GetAllFIOResident(ctx, "")
		if err != nil {
			c.log.Error("residentUsecase.GetAllFIOResident: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Список резидентов 💼")

		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = fioMarkup
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		return nil
	}
}

func (c *callbackResident) CallbackStartButton() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "<b>Список команд доступных для использования бота</b> ⏩")

		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = &handler.StartMenu
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		return nil
	}
}

func (c *callbackResident) CallbackStartChatGPT() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "💬 Нейросеть Avanti ответит на все ваши вопросы! Для начала работы напишите любой запрос...\n\n👌Если вы хотите остановить чат и воспользоваться другими командами используейте - /stop_chat_gpt")
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		return nil
	}
}

func (c *callbackResident) CallbackShowResidentByCluster() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		id := entity.FindID(update.CallbackData())
		if id == 0 {
			c.log.Error("entity.FindID == 0")
			handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrIncorrectCallbackData))
		}

		fioMarkup, err := c.residentUsecase.GetAllFIOResidentByCluster(ctx, "", id)
		if err != nil {
			c.log.Error("residentUsecase.GetAllFIOResident: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Список резидентов 💼")

		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = fioMarkup
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		return nil
	}
}

func (c *callbackResident) CallbackShowInstruction() tgbot.ViewFunc {
	return func(CallbackShowInstruction context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		text := "/start – список всех доступных команд\n\"Запустить Chat GPT\" – начать общение с искусственным интеллектом “AVANTI”\n\"Список кластеров\" – посмотреть список кластеров с резидентами бизнес клуба Avanti\n\"Показать расписание\" – посмотреть расписание мероприятий бизнес клуба “AVANTI”\n\"Список резидентов\" – посмотреть всех резидентов бизнес клуба “AVANTI”\n\"О нас\" – информация о бизнес клубе “AVANTI”\n\"Услуги AVANTI GROUP\" – список всех услуг предоставляемых бизнес клубом “AVANTI”\n\"Оставить заявку на вступление\" –  оставить заявку на вступление в бизнес клуб “AVANTI”\n\"Инструкция к боту\" – посмотреть список команд и описание к ним"

		_, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, text))
		if err != nil {
			return err
		}
		return nil
	}
}
