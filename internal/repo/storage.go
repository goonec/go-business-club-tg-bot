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
	GetAllByClusterID(ctx context.Context, id int) ([]entity.FIO, error)
}

type User interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetAllID(ctx context.Context) ([]int64, error)
}

type BusinessCluster interface {
	Delete(ctx context.Context, clusterID int) error
	Create(ctx context.Context, name string) (int, error)
	GetByName(ctx context.Context, name string) (*entity.BusinessCluster, error)
	GetAll(ctx context.Context) ([]entity.BusinessCluster, error)
}

type BusinessClusterResident interface {
	Create(ctx context.Context, IDBusinessCluster int, IDResident int) error
}

type Schedule interface {
	Create(ctx context.Context, file string) error
	Get(ctx context.Context) (*entity.Schedule, error)
}
