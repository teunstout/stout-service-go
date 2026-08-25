package usecase

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
)

type fakeAccountRepo struct {
	account domain.Account
	getErr  error
}

func (f *fakeAccountRepo) AccountExists(username string) (bool, error)          { return false, nil }
func (f *fakeAccountRepo) CreateAccount(username string, password string) error { return nil }
func (f *fakeAccountRepo) GetAccountByUsername(username string) (domain.Account, error) {
	if f.getErr != nil {
		return domain.Account{}, f.getErr
	}
	return f.account, nil
}

type fakeSessionRepo struct {
	createErr error
}

func (f *fakeSessionRepo) CreateSessionToken(token string, accountID int32, expiresAt time.Time) error {
	return f.createErr
}
func (f *fakeSessionRepo) GetAccountIdBySessionToken(token string) (int32, error) { return 0, nil }
func (f *fakeSessionRepo) DeleteSessionTokensByAccountId(accountID int32) error   { return nil }
func (f *fakeSessionRepo) DeleteSessionToken(token string) error                  { return nil }

type fakeCsrfRepo struct {
	createErr error
}

func (f *fakeCsrfRepo) CreateCsrfToken(token string, accountID int32, expiresAt time.Time) error {
	return f.createErr
}
func (f *fakeCsrfRepo) GetCsrfIdBySessionToken(token string) (int32, error) { return 0, nil }
func (f *fakeCsrfRepo) DeleteCsrfTokensByAccountId(accountID int32) error   { return nil }
func (f *fakeCsrfRepo) DeleteCsrfToken(token string) error                  { return nil }

func testSigningKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate test RSA key: %v", err)
	}
	return key
}

func TestLoginHappyPath(t *testing.T) {
	hash, err := domain.HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	accountRepo := &fakeAccountRepo{account: domain.Account{ID: 1, Username: "member", Password: hash}}
	u := NewLoginUsecase(accountRepo, &fakeSessionRepo{}, &fakeCsrfRepo{}, zap.NewNop(), testSigningKey(t))

	result, err := u.Login("member", "correct-password")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if result.Jwt == "" {
		t.Error("result.Jwt is empty")
	}
	if result.SessionToken == "" {
		t.Error("result.SessionToken is empty")
	}
	if result.CsrfToken == "" {
		t.Error("result.CsrfToken is empty")
	}
	if result.ExpiresAt.Before(time.Now()) {
		t.Error("result.ExpiresAt is in the past")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	hash, err := domain.HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	accountRepo := &fakeAccountRepo{account: domain.Account{ID: 1, Username: "member", Password: hash}}
	u := NewLoginUsecase(accountRepo, &fakeSessionRepo{}, &fakeCsrfRepo{}, zap.NewNop(), testSigningKey(t))

	if _, err := u.Login("member", "wrong-password"); err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestLoginUnknownUsername(t *testing.T) {
	accountRepo := &fakeAccountRepo{getErr: errors.New("no rows")}
	u := NewLoginUsecase(accountRepo, &fakeSessionRepo{}, &fakeCsrfRepo{}, zap.NewNop(), testSigningKey(t))

	if _, err := u.Login("nobody", "whatever"); err == nil {
		t.Fatal("expected error for unknown username, got nil")
	}
}

func TestLoginSessionRepoError(t *testing.T) {
	hash, err := domain.HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	accountRepo := &fakeAccountRepo{account: domain.Account{ID: 1, Username: "member", Password: hash}}
	sessionRepo := &fakeSessionRepo{createErr: errors.New("db down")}
	u := NewLoginUsecase(accountRepo, sessionRepo, &fakeCsrfRepo{}, zap.NewNop(), testSigningKey(t))

	if _, err := u.Login("member", "correct-password"); err == nil {
		t.Fatal("expected error when session repo fails, got nil")
	}
}
