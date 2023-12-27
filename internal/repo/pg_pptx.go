package repo

import (
	"context"
	"github.com/goonec/business-tg-bot/pkg/postgres"
)

type Pptx interface {
	Get(ctx context.Context) (string, error)
	Update(ctx context.Context, fileID string) error
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

func (p *pptxRepo) Update(ctx context.Context, fileID string) error {
	query := `update pptx set pptx_file_id = $1 where pptx_file_id = $2`

	file, err := p.Get(ctx)
	if err != nil {
		return err
	}

	_, err = p.Pool.Exec(ctx, query, fileID, file)
	return err
}
