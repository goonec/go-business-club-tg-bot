package usecase

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/repo"
)

type businessClusterUsecase struct {
	businessClusterRepo repo.BusinessCluster
}

func NewBusinessClusterUsecase(businessClusterRepo repo.BusinessCluster) BusinessCluster {
	return &businessClusterUsecase{
		businessClusterRepo: businessClusterRepo,
	}
}

func (b *businessClusterUsecase) Create(ctx context.Context, cluster string) error {
	_, err := b.businessClusterRepo.Create(ctx, cluster)
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

func (b *businessClusterUsecase) GetAllBusinessCluster(ctx context.Context) (*tgbotapi.InlineKeyboardMarkup, error) {
	businessCluster, err := b.businessClusterRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return b.createBusinessClusterMarkup(businessCluster)
}

func (b *businessClusterUsecase) createBusinessClusterMarkup(businessCluster []entity.BusinessCluster) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 3

	for i, el := range businessCluster {
		button := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", el.Name),
			fmt.Sprintf("cluster_%d", el.ID))

		row = append(row, button)

		if (i+1)%buttonsPerRow == 0 || i == len(businessCluster)-1 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	mainMenuButton := tgbotapi.NewInlineKeyboardButtonData("Вернуться к списку команд ⬆️", "main_menu")
	rows = append(rows, []tgbotapi.InlineKeyboardButton{mainMenuButton})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}
