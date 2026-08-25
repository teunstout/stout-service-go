package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"stout.dev/content/internal/domain"
)

type TranslationRepositoryInterface struct {
	pool *pgxpool.Pool
}

func NewTranslationRepository(pool *pgxpool.Pool) *TranslationRepositoryInterface {
	return &TranslationRepositoryInterface{pool: pool}
}

// BeginTx starts a transaction for the usecase layer to orchestrate a multi-statement
// sync (resolve list, then upsert each entry) atomically. The usecase layer never
// touches the pool directly - this is the only way in.
func (r *TranslationRepositoryInterface) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func UpdateTranslationEntry(
	ctx context.Context,
	tx pgx.Tx,
	accountID int32,
	listID int32,
	entryID int32,
	entry domain.TranslationEntryInput,
) (int64, error) {
	tag, err := tx.Exec(ctx, `
		UPDATE translation
		SET list_id = $1, original_html = $2, translation_html = $3, created_at = $4, updated_at = $5
		WHERE id = $6 AND list_id IN (SELECT id FROM translation_list WHERE account_id = $7)
	`, listID, entry.OriginalHTML, entry.TranslationHTML, entry.CreatedAt, entry.UpdatedAt, entryID, accountID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func InsertTranslationEntry(
	ctx context.Context,
	tx pgx.Tx,
	listID int32,
	entry domain.TranslationEntryInput,
) (int32, error) {
	var newID int32
	err := tx.QueryRow(ctx, `
		INSERT INTO translation (list_id, original_html, translation_html, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, listID, entry.OriginalHTML, entry.TranslationHTML, entry.CreatedAt, entry.UpdatedAt).Scan(&newID)
	return newID, err
}

func UpdateTranslationList(
	ctx context.Context,
	tx pgx.Tx,
	accountID int32,
	listID int32,
	name string,
) (int64, error) {
	tag, err := tx.Exec(ctx, `
		UPDATE translation_list SET name = $1 WHERE id = $2 AND account_id = $3
	`, name, listID, accountID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func InsertTranslationList(
	ctx context.Context,
	tx pgx.Tx,
	accountID int32,
	name string,
) (int32, error) {
	var newID int32
	err := tx.QueryRow(ctx, `
		INSERT INTO translation_list (account_id, name) VALUES ($1, $2) RETURNING id
	`, accountID, name).Scan(&newID)
	return newID, err
}

func (r *TranslationRepositoryInterface) GetLists(ctx context.Context, accountID int32) ([]domain.TranslationListOutput, error) {
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

	entryRows, err := r.pool.Query(ctx, `
		SELECT id, list_id, original_html, translation_html, created_at, updated_at
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
		if err := entryRows.Scan(
			&entry.ID, &listID, &entry.OriginalHTML, &entry.TranslationHTML, &entry.CreatedAt, &entry.UpdatedAt,
		); err != nil {
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

func (r *TranslationRepositoryInterface) DeleteList(ctx context.Context, accountID int32, listID int32) (bool, error) {
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

func (r *TranslationRepositoryInterface) DeleteEntries(ctx context.Context, accountID int32, ids []int32) ([]int32, error) {
	rows, err := r.pool.Query(ctx, `
		DELETE FROM translation
		WHERE id = ANY($1) AND list_id IN (SELECT id FROM translation_list WHERE account_id = $2)
		RETURNING id
	`, ids, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deletedIDs := []int32{}
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		deletedIDs = append(deletedIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deletedIDs, nil
}
