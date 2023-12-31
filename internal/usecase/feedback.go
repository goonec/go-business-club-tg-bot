package usecase

import (
	"context"
	"github.com/goonec/business-tg-bot/internal/entity"
	"github.com/goonec/business-tg-bot/internal/repo"
)

type feedbackUsecase struct {
	feedbackRepo repo.Feedback
}

func NewFeedbackUsecase(feedbackRepo repo.Feedback) Feedback {
	return &feedbackUsecase{
		feedbackRepo: feedbackRepo,
	}
}

func (f *feedbackUsecase) GetAllFeedback(ctx context.Context) ([]entity.Feedback, error) {
	return f.feedbackRepo.GetAll(ctx)
}

func (f *feedbackUsecase) DeleteFeedback(ctx context.Context, id int) error {
	return f.feedbackRepo.Delete(ctx, id)
}

func (f *feedbackUsecase) CreateFeedback(ctx context.Context, feedback *entity.Feedback) (*entity.Feedback, error) {
	return f.feedbackRepo.Create(ctx, feedback)
}
