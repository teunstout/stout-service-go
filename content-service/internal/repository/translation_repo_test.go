package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"stout.dev/content/internal/domain"
)

func setupTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	connString := os.Getenv("CONNECTION_STRING")
	if connString == "" {
		connString = "user=golang password=golang host=localhost port=5432 dbname=production sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Skipf("skipping: could not create pgx pool for %q: %v", connString, err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skipf("skipping: postgres not reachable at %q: %v", connString, err)
	}

	t.Cleanup(pool.Close)
	return pool
}

func createTestAccount(t *testing.T, pool *pgxpool.Pool) int32 {
	t.Helper()
	ctx := context.Background()

	username := fmt.Sprintf("test-translation-repo-%d", time.Now().UnixNano())
	var accountID int32
	err := pool.QueryRow(ctx, `
		INSERT INTO accounts (username, password) VALUES ($1, 'test-password') RETURNING id
	`, username).Scan(&accountID)
	if err != nil {
		t.Fatalf("failed to create test account: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx := context.Background()
		pool.Exec(cleanupCtx, `
			DELETE FROM translation WHERE list_id IN (SELECT id FROM translation_list WHERE account_id = $1)
		`, accountID)
		pool.Exec(cleanupCtx, `DELETE FROM translation_list WHERE account_id = $1`, accountID)
		pool.Exec(cleanupCtx, `DELETE FROM accounts WHERE id = $1`, accountID)
	})

	return accountID
}

func entryRow(t *testing.T, pool *pgxpool.Pool, id int32) (listID int32, originalHTML, translationHTML string, updatedAt time.Time) {
	t.Helper()
	err := pool.QueryRow(context.Background(), `
		SELECT list_id, original_html, translation_html, updated_at FROM translation WHERE id = $1
	`, id).Scan(&listID, &originalHTML, &translationHTML, &updatedAt)
	if err != nil {
		t.Fatalf("failed to read back translation row %d: %v", id, err)
	}
	return
}

func countEntriesForList(t *testing.T, pool *pgxpool.Pool, listID int32) int {
	t.Helper()
	var count int
	err := pool.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM translation WHERE list_id = $1
	`, listID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count entries for list %d: %v", listID, err)
	}
	return count
}

func TestSyncList_NewEntry_Inserts(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTranslationRepository(pool)
	accountID := createTestAccount(t, pool)
	ctx := context.Background()

	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := repo.SyncList(ctx, accountID, nil, "My List", []domain.TranslationEntryInput{
		{ID: nil, OriginalHTML: "hello", TranslationHTML: "world", CreatedAt: createdAt, UpdatedAt: createdAt},
	})
	if err != nil {
		t.Fatalf("SyncList returned error: %v", err)
	}
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry result, got %d", len(result.Entries))
	}
	if result.Entries[0].ID == 0 {
		t.Error("expected a non-zero assigned id for the new entry")
	}

	listID, originalHTML, translationHTML, _ := entryRow(t, pool, result.Entries[0].ID)
	if listID != result.ID {
		t.Errorf("expected entry's list_id to be %d, got %d", result.ID, listID)
	}
	if originalHTML != "hello" || translationHTML != "world" {
		t.Errorf("unexpected entry content: %q / %q", originalHTML, translationHTML)
	}
}

func TestSyncList_ExistingID_UpdatesContentAndUpdatedAt(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTranslationRepository(pool)
	accountID := createTestAccount(t, pool)
	ctx := context.Background()

	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	firstUpdatedAt := createdAt
	created, err := repo.SyncList(ctx, accountID, nil, "My List", []domain.TranslationEntryInput{
		{ID: nil, OriginalHTML: "original", TranslationHTML: "translation", CreatedAt: createdAt, UpdatedAt: firstUpdatedAt},
	})
	if err != nil {
		t.Fatalf("initial SyncList returned error: %v", err)
	}
	entryID := created.Entries[0].ID

	secondUpdatedAt := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	_, err = repo.SyncList(ctx, accountID, &created.ID, "My List", []domain.TranslationEntryInput{
		{ID: &entryID, OriginalHTML: "edited original", TranslationHTML: "edited translation", CreatedAt: createdAt, UpdatedAt: secondUpdatedAt},
	})
	if err != nil {
		t.Fatalf("second SyncList returned error: %v", err)
	}

	_, originalHTML, translationHTML, updatedAt := entryRow(t, pool, entryID)
	if originalHTML != "edited original" || translationHTML != "edited translation" {
		t.Errorf("expected content to be updated in place, got %q / %q", originalHTML, translationHTML)
	}
	if !updatedAt.Equal(secondUpdatedAt) {
		t.Errorf("expected updated_at to be %v, got %v", secondUpdatedAt, updatedAt)
	}
	if countEntriesForList(t, pool, created.ID) != 1 {
		t.Error("expected exactly 1 row for the list - an id-matched sync must update in place, not insert a second row")
	}
}

func TestSyncList_ExistingID_DifferentListID_MovesEntryWithoutDuplicating(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTranslationRepository(pool)
	accountID := createTestAccount(t, pool)
	ctx := context.Background()

	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	listA, err := repo.SyncList(ctx, accountID, nil, "List A", []domain.TranslationEntryInput{
		{ID: nil, OriginalHTML: "movable", TranslationHTML: "movable-t", CreatedAt: createdAt, UpdatedAt: createdAt},
	})
	if err != nil {
		t.Fatalf("SyncList (list A) returned error: %v", err)
	}
	entryID := listA.Entries[0].ID

	listB, err := repo.SyncList(ctx, accountID, nil, "List B", []domain.TranslationEntryInput{})
	if err != nil {
		t.Fatalf("SyncList (list B) returned error: %v", err)
	}

	_, err = repo.SyncList(ctx, accountID, &listB.ID, "List B", []domain.TranslationEntryInput{
		{ID: &entryID, OriginalHTML: "movable", TranslationHTML: "movable-t", CreatedAt: createdAt, UpdatedAt: createdAt},
	})
	if err != nil {
		t.Fatalf("SyncList (move) returned error: %v", err)
	}

	newListID, _, _, _ := entryRow(t, pool, entryID)
	if newListID != listB.ID {
		t.Errorf("expected entry to have moved to list B (%d), got list_id %d", listB.ID, newListID)
	}
	if countEntriesForList(t, pool, listA.ID) != 0 {
		t.Error("expected list A to have 0 entries after the move - a duplicate was left behind")
	}
	if countEntriesForList(t, pool, listB.ID) != 1 {
		t.Error("expected list B to have exactly 1 entry after the move")
	}
}

func TestSyncList_UnrecognizedID_FallsBackToInsert(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTranslationRepository(pool)
	accountID := createTestAccount(t, pool)
	ctx := context.Background()

	bogusID := int32(2_000_000_000)
	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := repo.SyncList(ctx, accountID, nil, "My List", []domain.TranslationEntryInput{
		{ID: &bogusID, OriginalHTML: "still inserted", TranslationHTML: "t", CreatedAt: createdAt, UpdatedAt: createdAt},
	})
	if err != nil {
		t.Fatalf("SyncList returned error: %v", err)
	}
	if result.Entries[0].ID == bogusID {
		t.Error("expected a freshly-assigned id, not the bogus one that was sent")
	}
	if countEntriesForList(t, pool, result.ID) != 1 {
		t.Error("expected the entry to have been inserted despite the unrecognized id")
	}
}

func TestSyncList_ForeignID_DoesNotOverwriteOtherAccountsEntry(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTranslationRepository(pool)
	victimAccountID := createTestAccount(t, pool)
	attackerAccountID := createTestAccount(t, pool)
	ctx := context.Background()

	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	victimList, err := repo.SyncList(ctx, victimAccountID, nil, "Victim's List", []domain.TranslationEntryInput{
		{ID: nil, OriginalHTML: "victim's content", TranslationHTML: "t", CreatedAt: createdAt, UpdatedAt: createdAt},
	})
	if err != nil {
		t.Fatalf("SyncList (victim) returned error: %v", err)
	}
	victimEntryID := victimList.Entries[0].ID

	attackerResult, err := repo.SyncList(ctx, attackerAccountID, nil, "Attacker's List", []domain.TranslationEntryInput{
		{ID: &victimEntryID, OriginalHTML: "attacker's content", TranslationHTML: "t", CreatedAt: createdAt, UpdatedAt: createdAt},
	})
	if err != nil {
		t.Fatalf("SyncList (attacker) returned error: %v", err)
	}

	if attackerResult.Entries[0].ID == victimEntryID {
		t.Fatal("attacker's entry was assigned the victim's exact id - ownership check did not fall back to insert")
	}

	victimListID, originalHTML, _, _ := entryRow(t, pool, victimEntryID)
	if victimListID != victimList.ID || originalHTML != "victim's content" {
		t.Error("victim's entry was overwritten by another account's sync")
	}
}

func TestSyncList_ResultEntriesMatchRequestOrder(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTranslationRepository(pool)
	accountID := createTestAccount(t, pool)
	ctx := context.Background()

	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := repo.SyncList(ctx, accountID, nil, "My List", []domain.TranslationEntryInput{
		{ID: nil, OriginalHTML: "first", TranslationHTML: "t1", CreatedAt: createdAt, UpdatedAt: createdAt},
		{ID: nil, OriginalHTML: "second", TranslationHTML: "t2", CreatedAt: createdAt, UpdatedAt: createdAt},
		{ID: nil, OriginalHTML: "third", TranslationHTML: "t3", CreatedAt: createdAt, UpdatedAt: createdAt},
	})
	if err != nil {
		t.Fatalf("SyncList returned error: %v", err)
	}
	if len(result.Entries) != 3 {
		t.Fatalf("expected 3 entry results, got %d", len(result.Entries))
	}

	wantContent := []string{"first", "second", "third"}
	for i, entryResult := range result.Entries {
		_, originalHTML, _, _ := entryRow(t, pool, entryResult.ID)
		if originalHTML != wantContent[i] {
			t.Errorf("result[%d] (id %d) has content %q, expected %q - result order doesn't match request order",
				i, entryResult.ID, originalHTML, wantContent[i])
		}
	}
}

func TestSyncList_ZeroEntries_ReturnsEmptyNotNilSlice(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTranslationRepository(pool)
	accountID := createTestAccount(t, pool)
	ctx := context.Background()

	result, err := repo.SyncList(ctx, accountID, nil, "Empty List", []domain.TranslationEntryInput{})
	if err != nil {
		t.Fatalf("SyncList returned error: %v", err)
	}
	if result.Entries == nil {
		t.Error("expected a non-nil empty slice (serializes to [], not JSON null)")
	}
	if len(result.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result.Entries))
	}
}

func TestGetLists_IncludesUpdatedAt(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTranslationRepository(pool)
	accountID := createTestAccount(t, pool)
	ctx := context.Background()

	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)
	_, err := repo.SyncList(ctx, accountID, nil, "My List", []domain.TranslationEntryInput{
		{ID: nil, OriginalHTML: "a", TranslationHTML: "b", CreatedAt: createdAt, UpdatedAt: updatedAt},
	})
	if err != nil {
		t.Fatalf("SyncList returned error: %v", err)
	}

	lists, err := repo.GetLists(ctx, accountID)
	if err != nil {
		t.Fatalf("GetLists returned error: %v", err)
	}
	if len(lists) != 1 || len(lists[0].Entries) != 1 {
		t.Fatalf("expected 1 list with 1 entry, got %d lists", len(lists))
	}
	if !lists[0].Entries[0].UpdatedAt.Equal(updatedAt) {
		t.Errorf("expected pulled entry's UpdatedAt to be %v, got %v", updatedAt, lists[0].Entries[0].UpdatedAt)
	}
}

func TestDeleteEntries_DeletesOnlyAccountOwnedIDsAndReturnsExactlyThose(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTranslationRepository(pool)
	accountID := createTestAccount(t, pool)
	otherAccountID := createTestAccount(t, pool)
	ctx := context.Background()

	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	mine, err := repo.SyncList(ctx, accountID, nil, "Mine", []domain.TranslationEntryInput{
		{ID: nil, OriginalHTML: "keep", TranslationHTML: "t", CreatedAt: createdAt, UpdatedAt: createdAt},
		{ID: nil, OriginalHTML: "delete-me", TranslationHTML: "t", CreatedAt: createdAt, UpdatedAt: createdAt},
	})
	if err != nil {
		t.Fatalf("SyncList (mine) returned error: %v", err)
	}
	theirs, err := repo.SyncList(ctx, otherAccountID, nil, "Theirs", []domain.TranslationEntryInput{
		{ID: nil, OriginalHTML: "not yours to delete", TranslationHTML: "t", CreatedAt: createdAt, UpdatedAt: createdAt},
	})
	if err != nil {
		t.Fatalf("SyncList (theirs) returned error: %v", err)
	}

	keepID := mine.Entries[0].ID
	deleteMeID := mine.Entries[1].ID
	foreignID := theirs.Entries[0].ID
	nonexistentID := int32(2_000_000_001)

	deletedIDs, err := repo.DeleteEntries(ctx, accountID, []int32{deleteMeID, foreignID, nonexistentID})
	if err != nil {
		t.Fatalf("DeleteEntries returned error: %v", err)
	}

	if len(deletedIDs) != 1 || deletedIDs[0] != deleteMeID {
		t.Errorf("expected deletedIDs to be exactly [%d], got %v", deleteMeID, deletedIDs)
	}
	if countEntriesForList(t, pool, mine.ID) != 1 {
		t.Error("expected exactly the owned entry to be deleted, leaving 1 behind")
	}
	if countEntriesForList(t, pool, theirs.ID) != 1 {
		t.Error("foreign account's entry must survive an unrelated account's delete request")
	}

	remainingListID, _, _, _ := entryRow(t, pool, keepID)
	if remainingListID != mine.ID {
		t.Error("the entry that should have been kept is gone or moved")
	}
}

func TestDeleteEntries_EmptyIDs_ReturnsEmptyNotNilSlice(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTranslationRepository(pool)
	accountID := createTestAccount(t, pool)

	deletedIDs, err := repo.DeleteEntries(context.Background(), accountID, []int32{})
	if err != nil {
		t.Fatalf("DeleteEntries returned error: %v", err)
	}
	if deletedIDs == nil {
		t.Error("expected a non-nil empty slice (serializes to [], not JSON null)")
	}
	if len(deletedIDs) != 0 {
		t.Errorf("expected 0 deleted ids, got %d", len(deletedIDs))
	}
}
