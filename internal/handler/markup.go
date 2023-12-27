package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	StartMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Запустить Chat GPT  🤖️", "chat_gpt")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список кластеров 🏆", "allcluster")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Показать раписание 🗓", "schedule")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список резидентов 💼", "resident")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Оставить заявку на вступление ✅", "request")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("О нас 💎", "pptx")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Услуги AVANTI GROUP 🛠️", "servicelist")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Инструкция к боту ❓", "instruction")))

	MainMenuButton = tgbotapi.NewInlineKeyboardButtonData("Вернуться к списку команд ⬆️", "main_menu")

	FeedbackButton = tgbotapi.NewInlineKeyboardButtonData("Обратная связь", "feedback")
)
