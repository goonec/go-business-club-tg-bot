package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type scheduleRepo struct {
	*postgres.Postgres
}

func NewScheduleRepo(pg *postgres.Postgres) Schedule {
	return &scheduleRepo{
		pg,
	}
}

func (s *scheduleRepo) Create(ctx context.Context, file string) error {
	query := `insert into schedule (photo_file_id) values ($1)`

	_, err := s.Pool.Exec(ctx, query, file)
	return err
}

func (s *scheduleRepo) Get(ctx context.Context) (*entity.Schedule, error) {
	query := `select id, photo_file_id, created_at 
				from schedule 
				order by created_at desc
				limit 1`
	var schedule entity.Schedule

	err := s.Pool.QueryRow(ctx, query).Scan(&schedule.ID, &schedule.PhotoFileID, &schedule.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, boterror.ErrNotFound
		}
		return nil, err
	}

	return &schedule, nil
}
