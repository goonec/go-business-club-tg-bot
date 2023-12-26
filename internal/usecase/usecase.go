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
	GetAllFIOResidentByCluster(ctx context.Context, command string, clusterID int) (*tgbotapi.InlineKeyboardMarkup, error)
}

type User interface {
	GetAllUserID(ctx context.Context) ([]int64, error)
	GetUser(ctx context.Context, id int64) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
}

type BusinessCluster interface {
	DeleteCluster(ctx context.Context, clusterID int) error
	CreateClusterResident(ctx context.Context, idBusinessCluster int, idResident int) error
	Create(ctx context.Context, cluster string) error
	GetAllBusinessCluster(ctx context.Context, callbackCommand string) (*tgbotapi.InlineKeyboardMarkup, error)
}

type Schedule interface {
	CreateSchedule(ctx context.Context, file string) error
	GetSchedule(ctx context.Context) (*entity.Schedule, error)
}

type Service interface {
	CreateService(ctx context.Context, name string) error
	CreateServiceDescribe(ctx context.Context, service *entity.ServiceDescribe) error
	DeleteService(ctx context.Context, id int) error
	DeleteServiceDescribe(ctx context.Context, id int) error
	GetAllService(ctx context.Context, command string) (*tgbotapi.InlineKeyboardMarkup, error)
	GetAllServiceDescribeByServiceID(ctx context.Context, serviceID int, command string) (*tgbotapi.InlineKeyboardMarkup, error)
	Get(ctx context.Context, serviceDescribeID int) (*entity.ServiceDescribe, error)
	GetAllServiceDescribe(ctx context.Context, command string) (*tgbotapi.InlineKeyboardMarkup, error)
}

type Feedback interface {
	GetAllFeedback(ctx context.Context) ([]entity.Feedback, error)
	DeleteFeedback(ctx context.Context, id int) error
	CreateFeedback(ctx context.Context, feedback *entity.Feedback) (*entity.Feedback, error)
}
