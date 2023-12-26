package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
)

func (b *Bot) middleware(update *tgbotapi.Update) error {
	channel := tgbotapi.ChatInfoConfig{
		tgbotapi.ChatConfig{
			ChatID: b.channelID,
		},
	}

	chat, err := b.api.GetChat(channel)
	if err != nil {
		b.log.Error("failed to get chat: %v", err)
		return err
	}

	cfg := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chat.ID,
			UserID: update.FromChat().ID,
		},
	}

	chatMember, err := b.api.GetChatMember(cfg)
	if err != nil {
		b.log.Error("error with chatID = %d:%v", chat.ID, err)
		return err
	}

	if chatMember.Status != "administrator" && chatMember.Status != "member" && chatMember.Status != "creator" {
		b.log.Error("[%s] %v", update.Message.From.UserName, boterror.ErrUserNotMember)
		return boterror.ErrUserNotMember
	}

	return nil
}
