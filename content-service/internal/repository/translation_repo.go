package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"stout.dev/content/internal/domain"
)

type TranslationRepository struct {
	pool *pgxpool.Pool
}

func NewTranslationRepository(pool *pgxpool.Pool) *TranslationRepository {
	return &TranslationRepository{pool: pool}
}

// SyncList resolves listID to a translation_list row owned by accountID (updating its name,
// or creating a new list if listID is nil or belongs to someone else/no one), then replaces
// all of that list's translation rows with entries. Runs as a single transaction so a
// mid-failure can't leave the list without any entries.
func (r *TranslationRepository) SyncList(
	ctx context.Context,
	accountID int32,
	listID *int32,
	name string,
	entries []domain.TranslationEntryInput,
) (domain.SyncListResult, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.SyncListResult{}, err
	}
	defer tx.Rollback(ctx)

	resolvedID, err := resolveListID(ctx, tx, accountID, listID, name)
	if err != nil {
		return domain.SyncListResult{}, err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM translation WHERE list_id = $1`, resolvedID); err != nil {
		return domain.SyncListResult{}, err
	}

	for _, entry := range entries {
		if _, err := tx.Exec(ctx, `
			INSERT INTO translation (list_id, original_html, translation_html)
			VALUES ($1, $2, $3)
		`, resolvedID, entry.OriginalHTML, entry.TranslationHTML); err != nil {
			return domain.SyncListResult{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.SyncListResult{}, err
	}

	return domain.SyncListResult{ID: resolvedID, Name: name, EntryCount: len(entries)}, nil
}

// resolveListID updates the name of the caller's own list if listID is set and matches an
// account_id-scoped row, otherwise inserts a new list. A stale or foreign listID collapses
// into the same "insert new" path as a nil listID (first sync) rather than erroring.
func resolveListID(
	ctx context.Context,
	tx pgx.Tx,
	accountID int32,
	listID *int32,
	name string,
) (int32, error) {
	if listID != nil {
		tag, err := tx.Exec(ctx, `
			UPDATE translation_list SET name = $1 WHERE id = $2 AND account_id = $3
		`, name, *listID, accountID)
		if err != nil {
			return 0, err
		}
		if tag.RowsAffected() == 1 {
			return *listID, nil
		}
	}

	var newID int32
	err := tx.QueryRow(ctx, `
		INSERT INTO translation_list (account_id, name) VALUES ($1, $2) RETURNING id
	`, accountID, name).Scan(&newID)
	return newID, err
}
