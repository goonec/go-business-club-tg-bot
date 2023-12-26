package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type serviceRepo struct {
	*postgres.Postgres
}

func NewServiceRepo(pg *postgres.Postgres) Service {
	return &serviceRepo{
		pg,
	}
}

func (s *serviceRepo) Delete(ctx context.Context, id int) error {
	query := `delete from service where id = $1`

	_, err := s.Pool.Exec(ctx, query, id)
	return err
}

func (s *serviceRepo) Create(ctx context.Context, service *entity.Service) error {
	query := `insert into service (name) values ($1)`

	_, err := s.Pool.Exec(ctx, query, service.Name)
	return err
}

func (s *serviceRepo) Get(ctx context.Context, id int) (*entity.Service, error) {
	query := `select id,name from service where id = $1`
	var service *entity.Service

	err := s.Pool.QueryRow(ctx, query, id).Scan(&service.ID, &service.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, boterror.ErrNotFound
		}
		return nil, err
	}

	return service, nil
}

func (s *serviceRepo) GetAll(ctx context.Context) ([]entity.Service, error) {
	query := `select id,name from service`

	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	service := make([]entity.Service, 0, 32)

	for rows.Next() {
		var s entity.Service

		err := rows.Scan(&s.ID, &s.Name)
		if err != nil {
			return nil, err
		}

		service = append(service, s)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return service, nil
}
