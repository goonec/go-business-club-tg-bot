package tgbot

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/openai"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error

type Bot struct {
	api          *tgbotapi.BotAPI
	log          *logger.Logger
	openAI       *openai.OpenAI
	cmdView      map[string]ViewFunc
	callbackView map[string]ViewFunc

	stateStore  map[int64]Store
	transportCh chan []string

	mu sync.RWMutex
}

type Store struct {
	waiting bool
	store   []string
}

func NewBot(api *tgbotapi.BotAPI, log *logger.Logger, openAI *openai.OpenAI, transportCh chan []string) *Bot {
	return &Bot{
		api:         api,
		openAI:      openAI,
		log:         log,
		transportCh: transportCh,
	}
}

func (b *Bot) RegisterCommandView(cmd string, view ViewFunc) {
	if b.cmdView == nil {
		b.cmdView = make(map[string]ViewFunc)
	}

	if b.stateStore == nil {
		b.stateStore = make(map[int64]Store)
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

			//b.mu.Lock()
			//delete(b.stateStore, update.FromChat().ID)
			//b.mu.Unlock()

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

	// Если пришло сообщение
	if update.Message != nil {
		b.log.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)

		nextStep := b.messageWithState(update)
		if !nextStep {
			return
		}

		// Провекрка на отсутствие команды и ожидания для запросов к openai, работает по аналагу default
		if _, ok := b.stateStore[update.Message.Chat.ID]; !ok {
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
		}

		var view ViewFunc

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
		// Если нажали на кнопку
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

func (b *Bot) messageWithState(update *tgbotapi.Update) bool {
	userID := update.Message.Chat.ID
	text := update.Message.Text

	if text == "/create_resident" {
		//b.mu.RLock()
		//defer b.mu.RUnlock()
		store, ok := b.stateStore[userID]
		if !ok {
			b.stateStore[userID] = Store{}
		}
		store.waiting = true

		return true
	}

	//b.mu.RLock()
	//defer b.mu.RUnlock()
	s, ok := b.stateStore[userID]
	if ok {
		switch {
		case len(s.store) == 0:
			s.store = append(s.store, text)

			msg := tgbotapi.NewMessage(userID, "[2] Введите резюме резидента.")
			if _, err := b.api.Send(msg); err != nil {
				b.log.Error("failed to send message: %v", err)
			}

			//b.mu.Lock()
			b.stateStore[userID] = s
			//b.mu.Unlock()

			return false
		case len(s.store) == 1:
			s.store = append(s.store, text)

			msg := tgbotapi.NewMessage(userID, "[3] Загрузите фотографию, связанную с резидентом.")
			if _, err := b.api.Send(msg); err != nil {
				b.log.Error("failed to send message: %v", err)
			}

			//b.mu.Lock()
			b.stateStore[userID] = s
			//b.mu.Unlock()

			return false
		case len(s.store) == 2:
			photo := update.Message.Photo
			if len(photo) > 0 {
				largestPhoto := photo[len(photo)-1]

				fileID := largestPhoto.FileID
				s.store = append(s.store, fileID)
			}
			s.waiting = false

			//b.mu.Lock()
			//defer b.mu.Unlock()
			b.stateStore[userID] = s
			b.transportCh <- b.stateStore[userID].store

			delete(b.stateStore, userID)
			return false
		}
	}
	return true
}
