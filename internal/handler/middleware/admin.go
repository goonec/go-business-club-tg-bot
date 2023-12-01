package middleware

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"

	"github.com/goonec/business-tg-bot/pkg/tgbot"
)

func AdminMiddleware(channelID int64, next tgbot.ViewFunc) tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		admins, err := bot.GetChatAdministrators(
			tgbotapi.ChatAdministratorsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: channelID,
				},
			})

		if err != nil {
			return err
		}

		for _, admin := range admins {
			if admin.User.ID == update.Message.From.ID {

				_, err = bot.Request(adminConfigMenu)
				if err != nil {
					log.Println("[ERROR] failed to request: %v", err)
				}

				return next(ctx, bot, update)
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID,
			"У вас нет прав на выполнение этой команды.")); err != nil {
			return err
		}

		return nil
	}
}

var adminConfigMenu = tgbotapi.NewSetMyCommands(
	tgbotapi.BotCommand{
		Command:     "/admin",
		Description: "Инструкция по использованию админки",
	},
	tgbotapi.BotCommand{
		Command:     "/chat_gpt",
		Description: "Общение с чатом",
	},
	//tgbotapi.BotCommand{
	//	Command:     "/create_resident_photo",
	//	Description: "Создание резедента только с фотографей",
	//},
	//tgbotapi.BotCommand{
	//	Command:     "/notify",
	//	Description: "Создание рассылки всем участникам бота",
	//},
	//tgbotapi.BotCommand{
	//	Command:     "/delete_resident",
	//	Description: "Удаление резидента",
	//},
	//tgbotapi.BotCommand{
	//	Command:     "/cancel",
	//	Description: "Используется в случае отмены администраторской команды",
	//},
)
