package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"stout.dev/content/internal/middleware"
	"stout.dev/content/internal/repository"
	"stout.dev/content/internal/usecase"
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

	username := fmt.Sprintf("test-translation-handler-%d", time.Now().UnixNano())
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

// bearerTokenFor signs a minimal JWT for accountID and returns the key pair's public half,
// so the returned token can be verified through the real middleware.AuthMiddleware exactly
// as it would be in production.
func bearerTokenFor(t *testing.T, accountID int32) (token string, verifyKey *rsa.PublicKey) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate test RSA key: %v", err)
	}

	claims := jwt.MapClaims{"sub": accountID}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		t.Fatalf("failed to sign test JWT: %v", err)
	}

	return signed, &privateKey.PublicKey
}

func TestHandleSyncList_UnrecognizedListID_Returns409(t *testing.T) {
	pool := setupTestPool(t)
	accountID := createTestAccount(t, pool)

	h := NewTranslationHandler(usecase.NewTranslationUsecase(repository.NewTranslationRepository(pool)))
	token, verifyKey := bearerTokenFor(t, accountID)
	wrapped := middleware.AuthMiddleware(h.HandleSyncList, verifyKey)

	bogusListID := int32(2_000_000_000)
	body, err := json.Marshal(map[string]any{
		"id":      bogusListID,
		"name":    "My List",
		"entries": []any{},
	})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/content/translation-lists/sync", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	wrapped(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusConflict, rec.Code, rec.Body.String())
	}

	var respBody map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to decode response body %q: %v", rec.Body.String(), err)
	}
	if respBody["error"] == "" {
		t.Error("expected a non-empty error message in the response body")
	}
}
