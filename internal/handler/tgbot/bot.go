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
	store *localstore.Store,
	log *logger.Logger,
	openAI *openai.OpenAI,
	userUsecase usecase.User,
	transportCh chan map[int64]map[string][]string,
	transportChResident chan map[int64]map[string][]string,
	transportChSchedule chan map[int64]map[string][]string,
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

		nextStep := b.messageWithState(update)
		if !nextStep {
			return
		}

		// Провекрка на отсутствие команды и ожидания для запросов к openai, работает по аналагу default
		if _, ok := b.readCommand(update.Message.Chat.ID, "chat_gpt"); ok {
			go func() {
				_, err = b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "✏️ Запрос создан, ожидайте…\n\n⏳ Среднее время ответа ChatGPT составляет от 7 до 19 секунд"))
				if err != nil {
					b.log.Error("failed to send message from ChatGPT %v", err)
					return
				}

				start := time.Now()
				openaiResponse, err := b.openAI.ResponseGPT(update.Message.Text)
				if err != nil {
					b.log.Error("failed to get response from GPT: %v", err)
				}

				_, err = b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, openaiResponse))
				if err != nil {
					b.log.Error("failed to send message from ChatGPT %v", err)
				}
				end := time.Since(start)
				b.log.Info("[%s] Время ответа: %f", update.Message.From.UserName, end.Seconds())
			}()
			return
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

func (b *Bot) messageWithState(update *tgbotapi.Update) bool {
	userID := update.Message.Chat.ID
	text := update.Message.Text

	isText := b.IsCommandText(text, userID)
	if isText != nil {
		return *isText
	}

	s, ok := b.read(userID)
	if ok {
		for key, value := range s {
			switch key {
			case "/create_schedule":
				photo := update.Message.Photo
				if len(photo) > 0 {
					largestPhoto := photo[len(photo)-1]

					fileID := largestPhoto.FileID
					b.set(fileID, key, userID)
				} else {
					b.delete(userID)

					msg := tgbotapi.NewMessage(userID, "Не является изображением [1].")
					if _, err := b.api.Send(msg); err != nil {
						b.log.Error("failed to send message: %v", err)
					}
					return false
				}

				d, _ := b.read(userID)
				b.transportChSchedule <- map[int64]map[string][]string{userID: d}

				b.delete(userID)
				return false
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

						msg := tgbotapi.NewMessage(userID, "Не является изображением [3].")
						if _, err := b.api.Send(msg); err != nil {
							b.log.Error("failed to send message: %v", err)
						}
						return false
					}

					d, _ := b.read(userID)
					b.transportChResident <- map[int64]map[string][]string{userID: d}

					b.delete(userID)
					return false
				}
			}

		}
		return true
	}
	return true
}
