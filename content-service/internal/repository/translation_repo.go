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
			INSERT INTO translation (list_id, original_html, translation_html, created_at)
			VALUES ($1, $2, $3, $4)
		`, resolvedID, entry.OriginalHTML, entry.TranslationHTML, entry.CreatedAt); err != nil {
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

// GetLists returns every list owned by accountID, each with its entries attached. Always
// returns a non-nil (possibly empty) slice so it serializes to `[]`, not JSON null.
func (r *TranslationRepository) GetLists(ctx context.Context, accountID int32) ([]domain.TranslationListOutput, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, created_at FROM translation_list WHERE account_id = $1 ORDER BY id
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lists := []domain.TranslationListOutput{}
	listIDs := []int32{}
	listIndexByID := map[int32]int{}
	for rows.Next() {
		var list domain.TranslationListOutput
		if err := rows.Scan(&list.ID, &list.Name, &list.CreatedAt); err != nil {
			return nil, err
		}
		list.Entries = []domain.TranslationEntryOutput{}
		listIndexByID[list.ID] = len(lists)
		lists = append(lists, list)
		listIDs = append(listIDs, list.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(listIDs) == 0 {
		return lists, nil
	}

	// list_id is only ever populated from listIDs above, which is already scoped to accountID.
	entryRows, err := r.pool.Query(ctx, `
		SELECT id, list_id, original_html, translation_html, created_at
		FROM translation
		WHERE list_id = ANY($1)
		ORDER BY created_at ASC, id ASC
	`, listIDs)
	if err != nil {
		return nil, err
	}
	defer entryRows.Close()

	for entryRows.Next() {
		var listID int32
		var entry domain.TranslationEntryOutput
		if err := entryRows.Scan(&entry.ID, &listID, &entry.OriginalHTML, &entry.TranslationHTML, &entry.CreatedAt); err != nil {
			return nil, err
		}
		idx := listIndexByID[listID]
		lists[idx].Entries = append(lists[idx].Entries, entry)
	}
	if err := entryRows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

// DeleteList removes a list and its entries if listID is owned by accountID. Returns
// false, nil (not an error) if no such list exists for this account — the caller may be
// retrying a delete that already succeeded, or racing another delete of the same list.
func (r *TranslationRepository) DeleteList(ctx context.Context, accountID int32, listID int32) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	var owned bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM translation_list WHERE id = $1 AND account_id = $2)
	`, listID, accountID).Scan(&owned)
	if err != nil {
		return false, err
	}
	if !owned {
		return false, nil
	}

	// translation.list_id has no ON DELETE CASCADE, so entries must go first.
	if _, err := tx.Exec(ctx, `DELETE FROM translation WHERE list_id = $1`, listID); err != nil {
		return false, err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM translation_list WHERE id = $1`, listID); err != nil {
		return false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}
