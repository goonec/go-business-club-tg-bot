package tgbot

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/usecase"
	"github.com/goonec/business-tg-bot/pkg/localstore"
	"github.com/goonec/business-tg-bot/pkg/logger"
	"github.com/goonec/business-tg-bot/pkg/openai"
	"runtime/debug"
	"sync"
	"time"
)

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error

var alwaysWork = map[string]struct{}{
	"/start":         struct{}{},
	"/exit":          struct{}{},
	"/resident_list": struct{}{},
}

type Bot struct {
	api          *tgbotapi.BotAPI
	log          *logger.Logger
	openAI       *openai.OpenAI
	userUsecase  usecase.User
	cmdView      map[string]ViewFunc
	callbackView map[string]ViewFunc

	store *localstore.Store

	channelID int64

	stateStore          map[int64]map[string][]string
	transportCh         chan map[int64]map[string][]string
	transportChResident chan map[int64]map[string][]string
	transportChSchedule chan map[int64]map[string][]string
	transportChFeedback chan map[int64][]interface{}
	transportPptx       chan map[int64]map[string][]string
	transportPhoto      chan map[int64]map[string][]string

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

func (b *Bot) deleteValue(userID int64, value string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	val, ok := b.stateStore[userID]
	if !ok {
		return
	}

	delete(val, value)
}

func NewBot(api *tgbotapi.BotAPI,
	store *localstore.Store,
	log *logger.Logger,
	openAI *openai.OpenAI,
	userUsecase usecase.User,
	transportCh chan map[int64]map[string][]string,
	transportChResident chan map[int64]map[string][]string,
	transportChSchedule chan map[int64]map[string][]string,
	transportChFeedback chan map[int64][]interface{},
	transportPptx chan map[int64]map[string][]string,
	transportPhoto chan map[int64]map[string][]string,
	channelID int64) *Bot {
	return &Bot{
		api:                 api,
		store:               store,
		log:                 log,
		openAI:              openAI,
		userUsecase:         userUsecase,
		transportCh:         transportCh,
		transportChResident: transportChResident,
		transportChSchedule: transportChSchedule,
		transportChFeedback: transportChFeedback,
		transportPptx:       transportPptx,
		transportPhoto:      transportPhoto,
		channelID:           channelID,
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

			b.log.Info("", b.stateStore)

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

	// Middleware для всех комманд
	//err := b.middleware(update)
	//if err != nil {
	//	return
	//}

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

		// Проверка на состояния админских команд у пользователя
		nextStepAdmin := b.messageWithState(update)
		if !nextStepAdmin {
			return
		}

		var executed bool
		for key, _ := range alwaysWork {
			// Проверка на команды, которые всегда должны использоваться без проверки пользовательских состояний
			if key == update.Message.Text {
				executed = true
				break
			}
		}

		if executed == false {
			nextStepUser := b.userMessageWithState(update)
			if !nextStepUser {
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
			if err == boterror.ErrIsNotAdmin {
				b.delete(update.Message.Chat.ID)
			}
			return
		}
		// Если нажали на кнопку
	} else if update.CallbackQuery != nil {
		b.log.Info("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackData())

		var callback ViewFunc

		err, callbackView := b.callbackHasString(update)
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
