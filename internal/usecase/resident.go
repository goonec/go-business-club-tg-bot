package usecase

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/repo"
)

type residentUsecase struct {
	residentRepo                repo.Resident
	businessClusterRepo         repo.BusinessCluster
	businessClusterResidentRepo repo.BusinessClusterResident
}

func NewResidentUsecase(residentRepo repo.Resident,
	businessClusterRepo repo.BusinessCluster,
	businessClusterResidentRepo repo.BusinessClusterResident) Resident {
	return &residentUsecase{
		residentRepo:                residentRepo,
		businessClusterRepo:         businessClusterRepo,
		businessClusterResidentRepo: businessClusterResidentRepo,
	}
}

func (r *residentUsecase) GetAllFIOResident(ctx context.Context, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	fio, err := r.residentRepo.GetAllFIO(ctx)
	if err != nil {
		return nil, err
	}

	return r.createFIOResidentMarkup(fio, command)
}

func (r *residentUsecase) GetAllFIOResidentByCluster(ctx context.Context, command string, clusterID int) (*tgbotapi.InlineKeyboardMarkup, error) {
	fio, err := r.residentRepo.GetAllByClusterID(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	return r.createFIOResidentMarkup(fio, command)
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
	residentID, err := r.residentRepo.Create(ctx, resident)
	if err != nil {
		errCode := repo.ErrorCode(err)
		if errCode == repo.ForeignKeyViolation {
			return boterror.ErrForeignKeyViolation
		}
		if errCode == repo.UniqueViolation {
			return boterror.ErrUniqueViolation
		}
		return err
	}

	bc := entity.BusinessCluster{}
	if resident.BusinessCluster == bc {
		fn := func(businessClusterID, residentID int) error {
			err := r.businessClusterResidentRepo.Create(context.Background(), businessClusterID, residentID)
			if err != nil {
				errCode := repo.ErrorCode(err)
				if errCode == repo.ForeignKeyViolation {
					return boterror.ErrForeignKeyViolation
				}
				if errCode == repo.UniqueViolation {
					return boterror.ErrUniqueViolation
				}
				return err
			}
			return nil
		}

		businessCluster, err := r.businessClusterRepo.GetByName(ctx, resident.BusinessCluster.Name)
		if err != nil {
			if errors.Is(err, boterror.ErrNotFound) {
				businessClusterID, err := r.businessClusterRepo.Create(ctx, resident.BusinessCluster.Name)
				if err != nil {
					errCode := repo.ErrorCode(err)
					if errCode == repo.ForeignKeyViolation {
						return boterror.ErrForeignKeyViolation
					}
					if errCode == repo.UniqueViolation {
						return boterror.ErrUniqueViolation
					}
					return err
				}
				err = fn(businessClusterID, residentID)
				if err != nil {
					return err
				}

				return nil
			}
			return err
		}

		err = fn(businessCluster.ID, residentID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *residentUsecase) DeleteResident(ctx context.Context, id int) error {
	err := r.residentRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *residentUsecase) createFIOResidentMarkup(fio []entity.FIO, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 3

	for i, el := range fio {
		button := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %s.", el.Firstname, el.Lastname),
			fmt.Sprintf("fio%s_%d", command, el.ID))

		row = append(row, button)

		if (i+1)%buttonsPerRow == 0 || i == len(fio)-1 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	mainMenuButton := tgbotapi.NewInlineKeyboardButtonData("Вернуться к списку команд ⬆️", "main_menu")
	rows = append(rows, []tgbotapi.InlineKeyboardButton{mainMenuButton})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}
