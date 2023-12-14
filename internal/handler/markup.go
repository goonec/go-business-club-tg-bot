package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	StartMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Запустить Chat GPT  🤖️", "chat_gpt")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список кластеров", "allcluster")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Показать раписание 🗓", "schedule")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список резидентов 💼", "resident")))

	MainMenuButton = tgbotapi.NewInlineKeyboardButtonData("Вернуться к списку команд ⬆️", "main_menu")
)
