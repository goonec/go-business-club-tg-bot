package usecase

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/repo"
)

type Pptx interface {
	GetPresentation(ctx context.Context) (string, error)
}

type pptxUsecase struct {
	pptxRepo repo.Pptx
}

func NewPPTXUsecase(pptxRepo repo.Pptx) Pptx {
	return &pptxUsecase{
		pptxRepo: pptxRepo,
	}
}

func (p *pptxUsecase) GetPresentation(ctx context.Context) (string, error) {
	return p.pptxRepo.Get(ctx)
}
