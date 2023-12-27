package tgbot

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (b *Bot) callbackHasString(update *tgbotapi.Update) (error, ViewFunc) {
	callbackData := update.CallbackData()

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
	case strings.HasPrefix(callbackData, "resident"):
		callbackView, ok := b.callbackView["resident"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "main_menu"):
		callbackView, ok := b.callbackView["main_menu"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "chat_gpt"):
		_, ok := b.read(update.CallbackQuery.Message.Chat.ID)
		if !ok {
			b.stateStore[update.CallbackQuery.Message.Chat.ID] = make(map[string][]string)
			b.stateStore[update.CallbackQuery.Message.Chat.ID]["chat_gpt"] = []string{}
		}
		callbackView, ok := b.callbackView["chat_gpt"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "feedback"):
		_, ok := b.read(update.CallbackQuery.Message.Chat.ID)
		if !ok {
			b.stateStore[update.CallbackQuery.Message.Chat.ID] = make(map[string][]string)
			b.stateStore[update.CallbackQuery.Message.Chat.ID]["feedback"] = []string{}
		}
		callbackView, ok := b.callbackView["feedback"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "request"):
		_, ok := b.read(update.CallbackQuery.Message.Chat.ID)
		if !ok {
			b.stateStore[update.CallbackQuery.Message.Chat.ID] = make(map[string][]string)
			b.stateStore[update.CallbackQuery.Message.Chat.ID]["request"] = []string{}
		}
		callbackView, ok := b.callbackView["request"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "allcluster"):
		callbackView, ok := b.callbackView["allcluster"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "cluster_"):
		callbackView, ok := b.callbackView["cluster"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "schedule"):
		callbackView, ok := b.callbackView["schedule"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "getcluster_"):
		callbackView, ok := b.callbackView["getcluster"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "fiogetresident_"):
		callbackView, ok := b.callbackView["fiogetresident"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "deletecluster_"):
		callbackView, ok := b.callbackView["deletecluster"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "servicelist"):
		callbackView, ok := b.callbackView["servicelist"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "service_"):
		callbackView, ok := b.callbackView["service"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "servicedescribelist_"):
		callbackView, ok := b.callbackView["servicedescribelist"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "servicedescribe_"):
		callbackView, ok := b.callbackView["servicedescribe"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "servicecreate_"):
		callbackView, ok := b.callbackView["servicecreate"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "pptx"):
		callbackView, ok := b.callbackView["pptx"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	case strings.HasPrefix(callbackData, "instruction"):
		callbackView, ok := b.callbackView["instruction"]
		if !ok {
			return errors.New("not found in map"), nil
		}
		return nil, callbackView
	}

	return nil, nil
}
