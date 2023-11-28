package usecase

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/entity"
)

type Resident interface {
	CreateResident(ctx context.Context, resident *entity.Resident) error
	GetResident(ctx context.Context, id int) (*entity.Resident, error)
	GetAllFIOResident(ctx context.Context) (*tgbotapi.InlineKeyboardMarkup, error)
}
