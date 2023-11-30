package tgbot

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/usecase"
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
	userUsecase  usecase.User
	cmdView      map[string]ViewFunc
	callbackView map[string]ViewFunc

	stateStore  map[int64]map[string][]string
	transportCh chan map[int64]map[string][]string

	mu sync.RWMutex
}

// set потокобезопасная запись структуры пользователя в map
func (b *Bot) set(data string, command string, userID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.stateStore[userID][command] = append(b.stateStore[userID][command], data)
}

// read потокобезопасные поиск пользователя в map
func (b *Bot) read(userID int64) (map[string][]string, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	data, ok := b.stateStore[userID]
	if !ok {
		return map[string][]string{}, false
	}

	return data, true
}

func (b *Bot) readCommand(userID int64, command string) ([]string, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	data, ok := b.stateStore[userID][command]
	if !ok {
		return []string{}, false
	}

	return data, true
}

// delete потокобезопасное удаление пользователя
func (b *Bot) delete(userID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.stateStore, userID)
}

func NewBot(api *tgbotapi.BotAPI,
	log *logger.Logger,
	openAI *openai.OpenAI,
	userUsecase usecase.User,
	transportCh chan map[int64]map[string][]string) *Bot {
	return &Bot{
		api:         api,
		log:         log,
		openAI:      openAI,
		userUsecase: userUsecase,
		transportCh: transportCh,
	}
}

func (b *Bot) RegisterCommandView(cmd string, view ViewFunc) {
	if b.cmdView == nil {
		b.cmdView = make(map[string]ViewFunc)
	}

	if b.stateStore == nil {
		b.stateStore = make(map[int64]map[string][]string)
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

	// Если пришло сообщение
	if update.Message != nil {
		b.log.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)

		_, err := b.userUsecase.GetUser(ctx, update.FromChat().ID)
		if err != nil {
			if errors.Is(err, boterror.ErrNotFound) {
				user := &entity.User{
					ID:         update.Message.From.ID,
					UsernameTG: update.Message.From.UserName,
				}

				err := b.userUsecase.CreateUser(ctx, user)
				if err != nil {
					b.log.Error("userUsecase.CreateUser: %v", err)
				}
			} else {
				b.log.Error("userUsecase.CreateUser: %v", err)
				return
			}
		}

		nextStep := b.messageWithState(update)
		if !nextStep {
			return
		}

		// Провекрка на отсутствие команды и ожидания для запросов к openai, работает по аналагу default
		if _, ok := b.read(update.Message.Chat.ID); !ok {
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
	case strings.HasPrefix(callbackData, "fio_"):
		callbackView, ok := b.callbackView["fio"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "fiodelete_"):
		callbackView, ok := b.callbackView["fiodelete"]
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

	if text == "/cancel" {
		b.cancelMessageWithState(userID)
		return false
	}

	if text == "/create_resident" {
		_, ok := b.read(userID)
		if !ok {
			b.stateStore[userID] = make(map[string][]string)
			b.stateStore[userID]["/create_resident"] = []string{}
		}
		return true
	}

	if text == "/create_resident_photo" {
		_, ok := b.read(userID)
		if !ok {
			b.stateStore[userID] = make(map[string][]string)
			b.stateStore[userID]["/create_resident_photo"] = []string{}
		}
		return true
	}

	if text == "/notify" {
		_, ok := b.read(userID)
		if !ok {
			b.stateStore[userID] = make(map[string][]string)
			b.stateStore[userID]["/notify"] = []string{}
		}
		return true
	}

	s, ok := b.read(userID)
	if ok {
		for key, value := range s {
			switch key {
			case "/notify":
				switch {
				case len(value) == 0:
					b.set(text, key, userID)

					msg := tgbotapi.NewMessage(userID, "[2] Загрузите фотографию, которая будет в рассылке.")
					if _, err := b.api.Send(msg); err != nil {
						b.log.Error("failed to send message: %v", err)
					}
					return false
				case len(value) == 1:
					photo := update.Message.Photo
					if len(photo) > 0 {
						largestPhoto := photo[len(photo)-1]

						fileID := largestPhoto.FileID
						b.set(fileID, key, userID)
					} else {
						b.delete(userID)

						msg := tgbotapi.NewMessage(userID, "Не является изображением [2]")
						if _, err := b.api.Send(msg); err != nil {
							b.log.Error("failed to send message: %v", err)
						}
						return false
					}

					d, _ := b.read(userID)
					b.log.Info("", d)
					b.transportCh <- map[int64]map[string][]string{userID: d}

					b.delete(userID)
					return false
				}
			case "/create_resident_photo":
				switch {
				case len(value) == 0:
					b.set(text, key, userID)

					msg := tgbotapi.NewMessage(userID, "[2] Загрузите фотографию, связанную с резидентом.")
					if _, err := b.api.Send(msg); err != nil {
						b.log.Error("failed to send message: %v", err)
					}
					return false
				case len(value) == 1:
					photo := update.Message.Photo
					if len(photo) > 0 {
						largestPhoto := photo[len(photo)-1]

						fileID := largestPhoto.FileID
						b.set(fileID, key, userID)
					} else {
						b.delete(userID)

						msg := tgbotapi.NewMessage(userID, "Не является изображением [2]")
						if _, err := b.api.Send(msg); err != nil {
							b.log.Error("failed to send message: %v", err)
						}
						return false
					}

					d, _ := b.read(userID)
					b.transportCh <- map[int64]map[string][]string{userID: d}

					b.delete(userID)
					return false
				}
			case "/create_resident":
				switch {
				case len(value) == 0:
					b.set(text, key, userID)

					msg := tgbotapi.NewMessage(userID, "[2] Введите резюме резидента.")
					if _, err := b.api.Send(msg); err != nil {
						b.log.Error("failed to send message: %v", err)
					}
					return false
				case len(value) == 1:
					b.set(text, key, userID)

					msg := tgbotapi.NewMessage(userID, "[3] Загрузите фотографию, связанную с резидентом.")
					if _, err := b.api.Send(msg); err != nil {
						b.log.Error("failed to send message: %v", err)
					}
					return false
				case len(value) == 2:
					photo := update.Message.Photo
					if len(photo) > 0 {
						largestPhoto := photo[len(photo)-1]

						fileID := largestPhoto.FileID
						b.set(fileID, key, userID)
					} else {
						b.delete(userID)

						msg := tgbotapi.NewMessage(userID, "Не является изображением [2].")
						if _, err := b.api.Send(msg); err != nil {
							b.log.Error("failed to send message: %v", err)
						}
						return false
					}

					d, _ := b.read(userID)
					b.transportCh <- map[int64]map[string][]string{userID: d}

					b.delete(userID)
					return false
				}
			}
		}
		return false
	}
	return true
}

func (b *Bot) cancelMessageWithState(userID int64) {
	b.delete(userID)

	msg := tgbotapi.NewMessage(userID, "Все команды отменены.")
	if _, err := b.api.Send(msg); err != nil {
		b.log.Error("failed to send message: %v", err)
	}
}
