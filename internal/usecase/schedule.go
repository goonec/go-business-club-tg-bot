package usecase

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/repo"
)

type scheduleUsecase struct {
	scheduleRepo repo.Schedule
}

func NewScheduleUsecase(scheduleRepo repo.Schedule) Schedule {
	return &scheduleUsecase{
		scheduleRepo: scheduleRepo,
	}
}

func (s *scheduleUsecase) CreateSchedule(ctx context.Context, file string) error {
	return s.scheduleRepo.Create(ctx, file)
}

func (s *scheduleUsecase) GetSchedule(ctx context.Context) (*entity.Schedule, error) {
	return s.scheduleRepo.Get(ctx)
}
