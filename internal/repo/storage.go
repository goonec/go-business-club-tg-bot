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
