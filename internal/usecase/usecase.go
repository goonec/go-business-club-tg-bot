package usecase

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/entity"
)

type Resident interface {
	CreateResident(ctx context.Context, resident *entity.Resident) error
	GetResident(ctx context.Context, id int) (*entity.Resident, error)
	GetAllFIOResident(ctx context.Context, command string) (*tgbotapi.InlineKeyboardMarkup, error)
	DeleteResident(ctx context.Context, id int) error
}

type User interface {
	GetAllUserID(ctx context.Context) ([]int64, error)
	GetUser(ctx context.Context, id int64) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
}

type BusinessCluster interface {
	GetAllBusinessCluster(ctx context.Context) (*tgbotapi.InlineKeyboardMarkup, error)
}
