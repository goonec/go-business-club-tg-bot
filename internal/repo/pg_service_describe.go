package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/pkg/postgres"
)

type serviceDescribeRepo struct {
	*postgres.Postgres
}

func NewServiceDescribeRepo(pg *postgres.Postgres) ServiceDescribe {
	return &serviceDescribeRepo{
		pg,
	}
}

func (s *serviceDescribeRepo) Create(ctx context.Context, service *entity.ServiceDescribe) error {
	query := `insert into service_describe (id_service,name,describe) values ($1,$2)`

	_, err := s.Pool.Exec(ctx, query, service.ServiceID, service.Describe)
	return err
}

func (s *serviceDescribeRepo) GetAllByServiceID(ctx context.Context, serviceID int) ([]entity.ServiceDescribe, error) {
	query := `select sd.id, sd.id_service,sd.describe,s.name, sd.name from service_describe sd
            	join service s on s.id = sd.id_service
            where sd.id_service = $1`

	rows, err := s.Pool.Query(ctx, query, serviceID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	service := make([]entity.ServiceDescribe, 0, 32)

	for rows.Next() {
		var s entity.ServiceDescribe

		err := rows.Scan(&s.ID, &s.ServiceID, &s.Describe, &s.Service.Name, &s.Name)
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

func (s *serviceDescribeRepo) Delete(ctx context.Context, id int) error {
	query := `delete from service_describe where id = $1`

	_, err := s.Pool.Exec(ctx, query, id)
	return err
}

func (s *serviceDescribeRepo) Get(ctx context.Context, id int) (*entity.ServiceDescribe, error) {
	query := `select sd.id, sd.id_service,sd.describe,s.name, sb.name from service_describe sd
            	join service s on s.id = sd.id_service
            where sd.id = $1`
	var sd *entity.ServiceDescribe

	err := s.Pool.QueryRow(ctx, query, id).Scan(&sd.ID, &sd.ServiceID, &sd.Describe, &sd.Service.Name, &sd.Name)
	return sd, err
}

func (s *serviceDescribeRepo) GetAll(ctx context.Context) ([]entity.ServiceDescribe, error) {
	query := `select id,name from service_describe`

	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	service := make([]entity.ServiceDescribe, 0, 64)

	for rows.Next() {
		var s entity.ServiceDescribe

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
