package usecase

import (
	"context"

	"stout.dev/content/internal/domain"
	"stout.dev/content/internal/repository"
)

type TranslationUsecase struct {
	repo *repository.TranslationRepository
}

func NewTranslationUsecase(repo *repository.TranslationRepository) *TranslationUsecase {
	return &TranslationUsecase{repo: repo}
}

func (u *TranslationUsecase) SyncList(
	ctx context.Context,
	accountID int32,
	req domain.SyncListRequest,
) (domain.SyncListResult, error) {
	return u.repo.SyncList(ctx, accountID, req.ID, req.Name, req.Entries)
}
