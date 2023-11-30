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
	userRepo         repo.User
	userResidentRepo repo.UserResident
}

func NewUserUsecase(userRepo repo.User, userResidentRepo repo.UserResident) User {
	return &userUsecase{
		userRepo:         userRepo,
		userResidentRepo: userResidentRepo,
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

	userResident, err := u.userResidentRepo.Get(ctx, &entity.UserResident{UserID: user.ID})
	if err != nil {
		if errors.Is(err, boterror.ErrNotFound) && userResident == nil {
			err := u.userResidentRepo.Create(ctx, user.ID)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if userResident.UserID != 0 {
		err := u.userResidentRepo.Update(ctx, user.ID)
		if err != nil {
			return err
		}
	}

	if userResident.UsernameTG != "" {
		err := u.userResidentRepo.Update(ctx, user.UsernameTG)
		if err != nil {
			return err
		}
	}

	return nil
}
