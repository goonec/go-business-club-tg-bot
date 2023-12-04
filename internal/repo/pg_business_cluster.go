package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/boterror"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type businessClusterRepo struct {
	*postgres.Postgres
}

func NewBusinessClusterRepository(postgres *postgres.Postgres) BusinessCluster {
	return &businessClusterRepo{
		postgres,
	}
}

func (b *businessClusterRepo) Create(ctx context.Context, name string) error {
	query := `insert into business_cluster (name) values ($1)`

	_, err := b.Pool.Exec(ctx, query, name)
	return err
}
func (b *businessClusterRepo) GetByName(ctx context.Context, name string) (*entity.BusinessCluster, error) {
	query := `select id, name from business_cluster where name = $1`
	var businessCluster entity.BusinessCluster

	err := b.Pool.QueryRow(ctx, query, name).Scan(&businessCluster.ID, &businessCluster.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, boterror.ErrNotFound
		}
		return nil, err
	}

	return &businessCluster, nil
}

func (b *businessClusterRepo) GetAll(ctx context.Context) ([]entity.BusinessCluster, error) {
	query := `select id, name from business_cluster`

	rows, err := b.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	businessClusters := make([]entity.BusinessCluster, 0, 32)

	for rows.Next() {
		var businessCluster entity.BusinessCluster

		err := rows.Scan(&businessCluster.ID, &businessCluster.Name)
		if err != nil {
			return nil, err
		}

		businessClusters = append(businessClusters, businessCluster)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return businessClusters, nil
}
