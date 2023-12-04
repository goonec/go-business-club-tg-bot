package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/pkg/postgres"
)

type businessClusterResidentRepo struct {
	*postgres.Postgres
}

func NewBusinessClusterResidentRepository(postgres *postgres.Postgres) BusinessClusterResident {
	return &businessClusterResidentRepo{
		postgres,
	}
}

func (b *businessClusterResidentRepo) Create(ctx context.Context, IDBusinessCluster int, IDResident int) error {
	query := `insert into business_cluster_resident (id_business_cluster,id_resident) values  ($1,$2)`

	_, err := b.Pool.Exec(ctx, query, IDBusinessCluster, IDResident)
	return err
}
