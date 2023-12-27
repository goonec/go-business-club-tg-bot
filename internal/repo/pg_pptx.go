package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/pkg/postgres"
)

type Pptx interface {
	Get(ctx context.Context) (string, error)
}

type pptxRepo struct {
	*postgres.Postgres
}

func NewPPTXRepo(pg *postgres.Postgres) Pptx {
	return &pptxRepo{
		pg,
	}
}

func (p *pptxRepo) Get(ctx context.Context) (string, error) {
	query := `select pptx_file_id from pptx`
	var pptx string

	err := p.Pool.QueryRow(ctx, query).Scan(&pptx)
	return pptx, err
}
