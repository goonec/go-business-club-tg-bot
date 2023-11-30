package usecase

import (
	"context"
	"errors"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/repo"
	"time"
)

type userUsecase struct {
	userRepo repo.User
}

func NewUserUsecase(userRepo repo.User) User {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetAllUserID(ctx context.Context) ([]int64, error) {
	allID, err := u.userRepo.GetAllID(ctx)
	if err != nil {
		return nil, err
	}

	return allID, nil
}

func (u *userUsecase) GetUser(ctx context.Context, id int64) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, boterror.ErrNotFound) {
			return nil, boterror.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) CreateUser(ctx context.Context, user *entity.User) error {
	user.CreatedAt = time.Now()

	err := u.userRepo.Create(ctx, user)
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
