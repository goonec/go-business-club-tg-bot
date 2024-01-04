package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (b *Bot) IsCommandText(text string, userID int64) *bool {
	switch text {
	//case "/create_service_photo":
	//	_, ok := b.read(userID)
	//	if !ok {
	//		b.stateStore[userID] = make(map[string][]string)
	//		b.stateStore[userID]["/create_service_photo"] = []string{}
	//	}
	//	return &[]bool{true}[0]
	case "/exit":
		b.cancelFeedbackOrRequest(userID)
		return &[]bool{false}[0]
	case "/cancel":
		b.cancelMessageWithState(userID)
		return &[]bool{false}[0]
	case "/stop_chat_gpt":
		b.cancelChatGptDialog(userID)
		return &[]bool{false}[0]
	case "/create_resident":
		_, ok := b.read(userID)
		if !ok {
			b.stateStore[userID] = make(map[string][]string)
			b.stateStore[userID]["/create_resident"] = []string{}
		}
		return &[]bool{true}[0]
	case "/create_resident_photo":
		_, ok := b.read(userID)
		if !ok {
			b.stateStore[userID] = make(map[string][]string)
			b.stateStore[userID]["/create_resident_photo"] = []string{}
		}
		return &[]bool{true}[0]
	case "/create_schedule":
		_, ok := b.read(userID)
		if !ok {
			b.stateStore[userID] = make(map[string][]string)
			b.stateStore[userID]["/create_schedule"] = []string{}
		}
		return &[]bool{true}[0]
	case "/notify":
		_, ok := b.read(userID)
		if !ok {
			b.stateStore[userID] = make(map[string][]string)
			b.stateStore[userID]["/notify"] = []string{}
		}
		return &[]bool{true}[0]
	case "/update_pptx":
		_, ok := b.read(userID)
		if !ok {
			b.stateStore[userID] = make(map[string][]string)
			b.stateStore[userID]["/update_pptx"] = []string{}
		}
		return &[]bool{true}[0]
	case "/create_under_service":
		_, ok := b.read(userID)
		if !ok {
			b.stateStore[userID] = make(map[string][]string)
			b.stateStore[userID]["/create_under_service"] = []string{}
		}
		return &[]bool{true}[0]
	//case "/chat_gpt":
	//	_, ok := b.read(userID)
	//	if !ok {
	//		b.stateStore[userID] = make(map[string][]string)
	//		b.stateStore[userID]["/chat_gpt"] = []string{}
	//	}
	//
	//	//msg := tgbotapi.NewMessage(userID, "Начните общение с Chat GPT!  💬\nЕсли вы хотите остановить чат и воспользоваться другими командами "+
	//	//	"используейте - /stop_chat_gpt")
	//	//if _, err := b.api.Send(msg); err != nil {
	//	//	b.log.Error("failed to send message: %v", err)
	//	//}
	//	return false
	default:
		return nil
	}
}

func (b *Bot) userMessageWithState(update *tgbotapi.Update) bool {
	// Провекрка на отсутствие команды и ожидания для запросов к openai, работает по аналагу default
	if _, ok := b.readCommand(update.Message.Chat.ID, "chat_gpt"); ok {
		go func(u *tgbotapi.Update) {
			sentMsg, err := b.api.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "✏️ Запрос создан, ожидайте…\n\n⏳ Среднее время ответа ChatGPT составляет от 7 до 19 секунд"))
			if err != nil {
				b.log.Error("failed to send message from ChatGPT %v", err)
				return
			}

			start := time.Now()
			openaiResponse, err := b.openAI.ResponseGPT(u.Message.Text)
			if err != nil {
				b.log.Error("failed to get response from GPT: %v", err)
			}

			updatedMsg := tgbotapi.EditMessageTextConfig{
				Text: openaiResponse,
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:    u.Message.Chat.ID,
					MessageID: sentMsg.MessageID,
				},
			}
			_, err = b.api.Send(updatedMsg)
			if err != nil {
				b.log.Error("failed to send message from ChatGPT %v", err)
			}

			end := time.Since(start)
			b.log.Info("[%s] Время ответа: %f", u.Message.From.UserName, end.Seconds())
		}(update)
		return false
	}

	// Обратная связь пользователей в кнопке услуги
	if _, ok := b.readCommand(update.Message.Chat.ID, "feedback"); ok {
		fb := []interface{}{update.Message.Text, *update, "услуги"}
		b.transportChFeedback <- map[int64][]interface{}{update.Message.Chat.ID: fb}
		b.delete(update.Message.Chat.ID)
		return false
	}

	// Заявка на вступление в клуб
	if _, ok := b.readCommand(update.Message.Chat.ID, "request"); ok {
		fb := []interface{}{update.Message.Text, *update, "заявка на вступление"}
		b.transportChFeedback <- map[int64][]interface{}{update.Message.Chat.ID: fb}
		b.delete(update.Message.Chat.ID)
		return false
	}

	return true
}

func (b *Bot) adminMessageWithState(update *tgbotapi.Update) bool {
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
			case "/create_under_service":
				switch {
				case len(value) == 0:
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

					msg := tgbotapi.NewMessage(userID, "[2] Загрузите информцию о услуги в формате, как представлено на примере.")
					if _, err := b.api.Send(msg); err != nil {
						b.log.Error("failed to send message: %v", err)
					}

					return false
				case len(value) == 1:
					b.set(text, key, userID)

					d, _ := b.read(userID)
					b.transportPhoto <- map[int64]map[string][]string{userID: d}

					b.delete(userID)
					return false
				}
			case "/update_pptx":
				pptx := update.Message.Document.FileID

				b.log.Info("", pptx)
				b.set(pptx, key, userID)

				d, _ := b.read(userID)
				b.transportPptx <- map[int64]map[string][]string{userID: d}

				b.delete(userID)
				return false
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

func (b *Bot) cancelMessageWithState(userID int64) {
	b.delete(userID)
	b.store.Delete(userID)

	msg := tgbotapi.NewMessage(userID, "Все команды отменены.")
	if _, err := b.api.Send(msg); err != nil {
		b.log.Error("failed to send message: %v", err)
	}
}

func (b *Bot) cancelChatGptDialog(userID int64) {
	b.delete(userID)

	msg := tgbotapi.NewMessage(userID, "Chat GPT остановлен.")
	if _, err := b.api.Send(msg); err != nil {
		b.log.Error("failed to send message: %v", err)
	}
}

func (b *Bot) cancelFeedbackOrRequest(userID int64) {
	b.delete(userID)

	msg := tgbotapi.NewMessage(userID, "Обратная связь отменена.")
	if _, err := b.api.Send(msg); err != nil {
		b.log.Error("failed to send message: %v", err)
	}
}
