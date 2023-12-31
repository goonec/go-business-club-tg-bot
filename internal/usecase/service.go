package usecase

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/repo"
)

type serviceUsecase struct {
	serviceRepo         repo.Service
	serviceDescribeRepo repo.ServiceDescribe
}

func NewServiceUsecase(serviceRepo repo.Service, serviceDescribeRepo repo.ServiceDescribe) Service {
	return &serviceUsecase{
		serviceRepo:         serviceRepo,
		serviceDescribeRepo: serviceDescribeRepo,
	}
}

func (s *serviceUsecase) CreateService(ctx context.Context, name string) error {
	err := s.serviceRepo.Create(ctx, name)
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

func (s *serviceUsecase) CreateServiceDescribe(ctx context.Context, service *entity.ServiceDescribe) error {
	err := s.serviceDescribeRepo.Create(ctx, service)
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

func (s *serviceUsecase) DeleteService(ctx context.Context, id int) error {
	return s.serviceRepo.Delete(ctx, id)
}

func (s *serviceUsecase) DeleteServiceDescribe(ctx context.Context, id int) error {
	return s.serviceDescribeRepo.Delete(ctx, id)
}

func (s *serviceUsecase) GetAllService(ctx context.Context, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	service, err := s.serviceRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return s.createServiceMarkup(service, command)
}

func (s *serviceUsecase) GetAllServiceDescribeByServiceID(ctx context.Context, serviceID int, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	serviceDescribe, err := s.serviceDescribeRepo.GetAllByServiceID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	return s.createServiceDescribeMarkup(serviceDescribe, command)
}

func (s *serviceUsecase) Get(ctx context.Context, serviceDescribeID int) (*entity.ServiceDescribe, error) {
	return s.serviceDescribeRepo.Get(ctx, serviceDescribeID)
}

func (s *serviceUsecase) CreatePhoto(ctx context.Context, id int, fileID string) error {
	return s.serviceDescribeRepo.CreatePhoto(ctx, id, fileID)
}

func (s *serviceUsecase) GetAllServiceDescribe(ctx context.Context, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	serviceDescribe, err := s.serviceDescribeRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return s.createServiceDescribeMarkup(serviceDescribe, command)
}

func (s *serviceUsecase) createServiceMarkup(service []entity.Service, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1

	for i, el := range service {
		button := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", el.Name),
			fmt.Sprintf("service%s_%d", command, el.ID))

		row = append(row, button)

		if (i+1)%buttonsPerRow == 0 || i == len(service)-1 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if command != "create" {
		mainMenuButton := tgbotapi.NewInlineKeyboardButtonData("Вернуться к списку команд ⬆️", "main_menu")
		rows = append(rows, []tgbotapi.InlineKeyboardButton{mainMenuButton})
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}

func (s *serviceUsecase) createServiceDescribeMarkup(service []entity.ServiceDescribe, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1

	for i, el := range service {
		button := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", el.Name),
			fmt.Sprintf("service%s_%d", command, el.ID))

		row = append(row, button)

		if (i+1)%buttonsPerRow == 0 || i == len(service)-1 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	mainMenuButton := tgbotapi.NewInlineKeyboardButtonData("Вернуться к списку команд ⬆️", "main_menu")
	rows = append(rows, []tgbotapi.InlineKeyboardButton{mainMenuButton})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}
