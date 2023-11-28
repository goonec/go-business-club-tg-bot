package tgbot

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/openai"
	"runtime/debug"
	"strings"
	"time"
)

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error

type Bot struct {
	api          *tgbotapi.BotAPI
	log          *logger.Logger
	openAI       *openai.OpenAI
	cmdView      map[string]ViewFunc
	callbackView map[string]ViewFunc
}

func NewBot(api *tgbotapi.BotAPI, log *logger.Logger, openAI *openai.OpenAI) *Bot {
	return &Bot{
		api:    api,
		openAI: openAI,
		log:    log,
	}
}

func (b *Bot) RegisterCommandView(cmd string, view ViewFunc) {
	if b.cmdView == nil {
		b.cmdView = make(map[string]ViewFunc)
	}

	b.cmdView[cmd] = view
}

func (b *Bot) RegisterCommandCallback(callback string, view ViewFunc) {
	if b.callbackView == nil {
		b.callbackView = make(map[string]ViewFunc)
	}

	b.callbackView[callback] = view
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			b.handlerUpdate(updateCtx, &update)
			cancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) handlerUpdate(ctx context.Context, update *tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			b.log.Error("panic recovered: %v, %s", p, string(debug.Stack()))
		}
	}()

	if update.Message != nil {
		b.log.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)

		var view ViewFunc

		if !update.Message.IsCommand() {
			openaiResponse, err := b.openAI.ResponseGPT(update.Message.Text)
			if err != nil {
				b.log.Error("failed to get response from GPT: %v", err)
			}

			_, err = b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, openaiResponse))
			if err != nil {
				b.log.Error("failed to send message from ChatGPT %v", err)
			}
			return
		}

		cmd := update.Message.Command()

		cmdView, ok := b.cmdView[cmd]
		if !ok {
			return
		}

		view = cmdView

		if err := view(ctx, b.api, update); err != nil {
			b.log.Error("failed to handle update: %v", err)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "internal error")
			if _, err := b.api.Send(msg); err != nil {
				b.log.Error("failed to send message: %v", err)
			}
			return
		}
	} else if update.CallbackQuery != nil {
		b.log.Info("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackData())

		var callback ViewFunc
		callbackData := update.CallbackData()

		err, callbackView := b.callbackHasString(callbackData)
		if err != nil {
			b.log.Error("%v", err)
			return
		}

		callback = callbackView

		if err := callback(ctx, b.api, update); err != nil {
			b.log.Error("failed to handle update: %v", err)

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "internal error")
			if _, err := b.api.Send(msg); err != nil {
				b.log.Error("failed to send message: %v", err)
			}
			return
		}

	}

}

func (b *Bot) callbackHasString(callbackData string) (error, ViewFunc) {
	switch {
	case strings.HasPrefix(callbackData, "fio"):
		callbackView, ok := b.callbackView["fio"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	}

	return nil, nil
}
