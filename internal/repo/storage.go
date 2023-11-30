package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/entity"
)

type Resident interface {
	GetAll(ctx context.Context) ([]entity.Resident, error)
	GetAllFIO(ctx context.Context) ([]entity.FIO, error)
	GetByID(ctx context.Context, id int) (*entity.Resident, error)
	Create(ctx context.Context, resident *entity.Resident) error
}

type User interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetAllID(ctx context.Context) ([]int64, error)
}

type UserResident interface {
	Create(ctx context.Context, usernameTG string, userID int64) error
	CreateWithUsername(ctx context.Context, usernameTG string) error
	Update(ctx context.Context, data any) error
	Get(ctx context.Context, userResident *entity.UserResident) (*entity.UserResident, error)
}
