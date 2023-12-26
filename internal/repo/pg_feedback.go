package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type feedbackRepo struct {
	*postgres.Postgres
}

func NewFeedbackRepo(pg *postgres.Postgres) Feedback {
	return &feedbackRepo{
		pg,
	}
}

func (f *feedbackRepo) GetAll(ctx context.Context) ([]entity.Feedback, error) {
	query := `select id,message,created_at,tg_username,type from feedback`

	rows, err := f.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	fio := make([]entity.Feedback, 0, 128)

	for rows.Next() {
		var f entity.Feedback

		err := rows.Scan(&f.ID,
			&f.Message,
			&f.CreatedAt,
			&f.UsernameTG,
			&f.Type,
		)
		if err != nil {
			return nil, err
		}

		fio = append(fio, f)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return fio, nil

}

func (f *feedbackRepo) Delete(ctx context.Context, id int) error {
	query := `delete from feedback where id = $1`

	_, err := f.Pool.Exec(ctx, query, id)
	return err
}

func (f *feedbackRepo) Create(ctx context.Context, feedback *entity.Feedback) (*entity.Feedback, error) {
	query := `insert into feedback (message,tg_username,type) values ($1,$2,$3) returning *`

	var fb entity.Feedback

	err := f.Pool.QueryRow(ctx, query, feedback.Message, feedback.UsernameTG, feedback.Type).Scan(
		&fb.ID,
		&fb.Message,
		&fb.Type,
		&fb.CreatedAt,
		&fb.UsernameTG)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, boterror.ErrNotFound
		}
		return nil, err
	}

	return &fb, nil
}
