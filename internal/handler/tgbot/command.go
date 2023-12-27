package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) IsCommandText(text string, userID int64) *bool {
	switch text {
	//case "/create_service_photo":
	//	_, ok := b.read(userID)
	//	if !ok {
	//		b.stateStore[userID] = make(map[string][]string)
	//		b.stateStore[userID]["/create_service_photo"] = []string{}
	//	}
	//	return &[]bool{true}[0]
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
	case "/create_pptx":
		_, ok := b.read(userID)
		if !ok {
			b.stateStore[userID] = make(map[string][]string)
			b.stateStore[userID]["/create_pptx"] = []string{}
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
