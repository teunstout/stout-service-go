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

func (u *TranslationUsecase) GetLists(ctx context.Context, accountID int32) (domain.GetListsResult, error) {
	lists, err := u.repo.GetLists(ctx, accountID)
	if err != nil {
		return domain.GetListsResult{}, err
	}
	return domain.GetListsResult{Lists: lists}, nil
}

func (u *TranslationUsecase) DeleteList(ctx context.Context, accountID int32, listID int32) (bool, error) {
	return u.repo.DeleteList(ctx, accountID, listID)
}

func (u *TranslationUsecase) DeleteEntries(ctx context.Context, accountID int32, ids []int32) ([]int32, error) {
	return u.repo.DeleteEntries(ctx, accountID, ids)
}
