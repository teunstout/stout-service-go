package usecase

import (
	"context"

	"github.com/jackc/pgx/v5"
	"stout.dev/content/internal/domain"
	"stout.dev/content/internal/repository"
)

type TranslationUsecase struct {
	repo repository.TranslationRepository
}

func NewTranslationUsecase(repo repository.TranslationRepository) *TranslationUsecase {
	return &TranslationUsecase{repo: repo}
}

func (u *TranslationUsecase) SyncList(
	ctx context.Context,
	accountID int32,
	req domain.SyncListRequest,
) (domain.SyncListResult, error) {
	tx, err := u.repo.BeginTx(ctx)
	if err != nil {
		return domain.SyncListResult{}, err
	}
	defer tx.Rollback(ctx)

	resolvedID, err := resolveListID(ctx, tx, accountID, req.ID, req.Name)
	if err != nil {
		return domain.SyncListResult{}, err
	}

	results := make([]domain.SyncEntryResult, len(req.Entries))
	for i, entry := range req.Entries {
		entryID, err := upsertEntry(ctx, tx, accountID, resolvedID, entry)
		if err != nil {
			return domain.SyncListResult{}, err
		}
		results[i] = domain.SyncEntryResult{ID: entryID}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.SyncListResult{}, err
	}

	return domain.SyncListResult{ID: resolvedID, Name: req.Name, Entries: results}, nil
}

// upsertEntry holds the business decision of whether a synced entry is an update to an
// existing row or the creation of a new one. Repository functions below only ever do the
// one database action they're named for.
func upsertEntry(
	ctx context.Context,
	tx pgx.Tx,
	accountID int32,
	listID int32,
	entry domain.TranslationEntryInput,
) (int32, error) {
	if entry.ID != nil {
		rowsAffected, err := repository.UpdateTranslationEntry(ctx, tx, accountID, listID, *entry.ID, entry)
		if err != nil {
			return 0, err
		}
		if rowsAffected == 1 {
			return *entry.ID, nil
		}
		return 0, domain.ErrEntryNotFound
	}

	return repository.InsertTranslationEntry(ctx, tx, listID, entry)
}

func resolveListID(
	ctx context.Context,
	tx pgx.Tx,
	accountID int32,
	listID *int32,
	name string,
) (int32, error) {
	if listID != nil {
		rowsAffected, err := repository.UpdateTranslationList(ctx, tx, accountID, *listID, name)
		if err != nil {
			return 0, err
		}
		if rowsAffected == 1 {
			return *listID, nil
		}
		return 0, domain.ErrListNotFound
	}

	return repository.InsertTranslationList(ctx, tx, accountID, name)
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
