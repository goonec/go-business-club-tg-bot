package middleware

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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

		fmt.Println(err)
		if err != nil {
			return err
		}

		for _, admin := range admins {
			if admin.User.ID == update.Message.From.ID {
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
