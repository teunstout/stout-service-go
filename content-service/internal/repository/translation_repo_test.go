package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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
