package usecase

import (
	"context"
	"errors"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/repo"
)

type residentUsecase struct {
	residentRepo repo.Resident
}

func NewResidentUsecase(residentRepo repo.Resident) Resident {
	return &residentUsecase{
		residentRepo: residentRepo,
	}
}

func (r *residentUsecase) GetAllFIOResident(ctx context.Context) ([]entity.FIO, error) {
	fio, err := r.residentRepo.GetAllFIO(ctx)
	if err != nil {
		return nil, err
	}

	return fio, nil
}

func (r *residentUsecase) GetResident(ctx context.Context, id int) (*entity.Resident, error) {
	resident, err := r.residentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, boterror.ErrNotFound) {
			return nil, boterror.ErrNotFound
		}
		return nil, err
	}

	return resident, nil
}

func (r *residentUsecase) CreateResident(ctx context.Context, resident *entity.Resident) error {
	err := r.residentRepo.Create(ctx, resident)
	if err != nil {
		return err
	}

	return nil
}
