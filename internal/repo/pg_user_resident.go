package repo

import (
	"context"
	"database/sql"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/pkg/postgres"
)

type userResidentRepository struct {
	*postgres.Postgres
}

func NewUserResident(postgres *postgres.Postgres) UserResident {
	return &userResidentRepository{
		postgres,
	}
}

func (u *userResidentRepository) Create(ctx context.Context, data any) error {
	queryUserID := `insert into user_resident (user_id) values ($1)`
	queryUsernameTG := `insert into user_resident (user_id) values ($1)`

	switch data {
	case data.(int64):
		_, err := u.Pool.Exec(ctx, queryUserID, data.(int64))
		return err
	case data.(string):
		_, err := u.Pool.Exec(ctx, queryUsernameTG, data.(string))
		return err
	}

	return nil
}

func (u *userResidentRepository) Update(ctx context.Context, data any) error {
	queryUserID := `update user_resident SET user_id = $1 where id = $2`
	queryUsernameTG := `update user_resident SET tg_username = $1 where id = $2`

	switch data {
	case data.(int64):
		_, err := u.Pool.Exec(ctx, queryUserID, data.(int64))
		return err
	case data.(string):
		_, err := u.Pool.Exec(ctx, queryUsernameTG, data.(string))
		return err
	}

	return nil
}

func (u *userResidentRepository) Get(ctx context.Context, userResident *entity.UserResident) (*entity.UserResident, error) {
	query := `select id, user_id, tg_username from user_resident where user_id = $1 or tg_username = $2;`
	var ur *entity.UserResident

	err := u.Pool.QueryRow(ctx, query,
		userResident.UserID,
		userResident.UsernameTG).Scan(&ur.ID, &ur.UserID, &ur.UsernameTG)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, boterror.ErrNotFound
		}
		return nil, err
	}

	return ur, nil
}
