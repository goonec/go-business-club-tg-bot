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
	residentRepo repo.Resident
}

func NewResidentUsecase(residentRepo repo.Resident) Resident {
	return &residentUsecase{
		residentRepo: residentRepo,
	}
}

func (r *residentUsecase) GetAllFIOResident(ctx context.Context) (*tgbotapi.InlineKeyboardMarkup, error) {
	fio, err := r.residentRepo.GetAllFIO(ctx)
	if err != nil {
		return nil, err
	}

	return r.createFIOResidentMarkup(fio)
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

func (r *residentUsecase) createFIOResidentMarkup(fio []entity.FIO) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for _, el := range fio {
		button := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %s.%s.", el.Firstname, el.Lastname, el.Patronymic), fmt.Sprintf("fio_%d", el.ID))
		row = append(row, button)
		rows = append(rows, row)
		row = []tgbotapi.InlineKeyboardButton{}
	}

	rows = append(rows, row)

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}
