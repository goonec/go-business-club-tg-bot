package repo

import (
	"context"
	"database/sql"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/pkg/postgres"
)

type userRepository struct {
	*postgres.Postgres
}

func NewUserRepository(postgres *postgres.Postgres) User {
	return &userRepository{
		postgres,
	}
}

func (u *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `insert into "user" (id, tg_username, create_at, role) values ($1,$2,$3,$4)`

	_, err := u.Pool.Exec(ctx, query, user.ID, user.UsernameTG, user.CreatedAt, user.Role)

	return err
}

func (u *userRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `select id, tg_username, create_at, role from "user" where id = $1`
	var user entity.User

	err := u.Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.UsernameTG, &user.CreatedAt, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, boterror.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) GetAllID(ctx context.Context) ([]int64, error) {
	query := `select id from "user"`

	rows, err := u.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allID := make([]int64, 0, 256)

	for rows.Next() {
		var id int64

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		allID = append(allID, id)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return allID, nil
}
