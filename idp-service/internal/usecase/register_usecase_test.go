package usecase

import (
	"errors"
	"testing"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
)

type fakeRegisterAccountRepo struct {
	exists      bool
	existsErr   error
	createErr   error
	created     bool
	createdUser string
}

func (f *fakeRegisterAccountRepo) AccountExists(username string) (bool, error) {
	return f.exists, f.existsErr
}
func (f *fakeRegisterAccountRepo) CreateAccount(username string, password string) error {
	f.created = true
	f.createdUser = username
	return f.createErr
}
func (f *fakeRegisterAccountRepo) GetAccountByUsername(username string) (domain.Account, error) {
	return domain.Account{}, nil
}

func TestRegisterHappyPath(t *testing.T) {
	repo := &fakeRegisterAccountRepo{}
	u := NewRegisterUsecase(repo, zap.NewNop(), testSigningKey(t))

	if err := u.Register("newuser", "s3cret-password"); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if !repo.created {
		t.Error("expected CreateAccount to be called")
	}
	if repo.createdUser != "newuser" {
		t.Errorf("CreateAccount called with username %q, want %q", repo.createdUser, "newuser")
	}
}

func TestRegisterAccountAlreadyExists(t *testing.T) {
	repo := &fakeRegisterAccountRepo{exists: true}
	u := NewRegisterUsecase(repo, zap.NewNop(), testSigningKey(t))

	err := u.Register("newuser", "s3cret-password")
	if !errors.Is(err, domain.ErrAccountAlreadyExists) {
		t.Fatalf("Register error = %v, want domain.ErrAccountAlreadyExists", err)
	}
	if repo.created {
		t.Error("CreateAccount should not be called when account already exists")
	}
}

func TestRegisterAccountExistsCheckError(t *testing.T) {
	repo := &fakeRegisterAccountRepo{existsErr: errors.New("db down")}
	u := NewRegisterUsecase(repo, zap.NewNop(), testSigningKey(t))

	if err := u.Register("newuser", "s3cret-password"); err == nil {
		t.Fatal("expected error when AccountExists fails, got nil")
	}
}

func TestRegisterCreateAccountError(t *testing.T) {
	repo := &fakeRegisterAccountRepo{createErr: errors.New("db down")}
	u := NewRegisterUsecase(repo, zap.NewNop(), testSigningKey(t))

	if err := u.Register("newuser", "s3cret-password"); err == nil {
		t.Fatal("expected error when CreateAccount fails, got nil")
	}
}
