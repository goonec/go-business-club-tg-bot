package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/entity"
)

type Resident interface {
	GetAll(ctx context.Context) ([]entity.Resident, error)
	GetAllFIO(ctx context.Context) ([]entity.FIO, error)
	GetByID(ctx context.Context, id int) (*entity.Resident, error)
	Create(ctx context.Context, resident *entity.Resident) (int, error)
	DeleteByID(ctx context.Context, id int) error
}

type User interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetAllID(ctx context.Context) ([]int64, error)
}

type BusinessCluster interface {
	Create(ctx context.Context, name string) error
	GetByName(ctx context.Context, name string) (*entity.BusinessCluster, error)
	GetAll(ctx context.Context) ([]entity.BusinessCluster, error)
}

type BusinessClusterResident interface {
}
