package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
)

type AccountRepositoryInterface struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewAccountRepository(pool *pgxpool.Pool, l *zap.Logger) *AccountRepositoryInterface {
	return &AccountRepositoryInterface{
		pool:   pool,
		logger: l,
	}
}

func (r *AccountRepositoryInterface) AccountExists(username string) (bool, error) {
	r.logger.Debug("Checking if account exists", zap.String("username", username))
	var exists bool
	err := r.pool.QueryRow(context.Background(), `
		SELECT EXISTS(SELECT 1 FROM accounts WHERE username = LOWER($1))
	`, username).Scan(&exists)
	return exists, err
}

func (r *AccountRepositoryInterface) CreateAccount(username string, password string) error {
	r.logger.Debug("Creating account", zap.String("username", username))
	_, err := r.pool.Exec(context.Background(), `
		INSERT INTO accounts (username, password)
		VALUES (LOWER($1), $2)
	`, username, password)
	return err
}

func (r *AccountRepositoryInterface) GetAccountByUsername(username string) (domain.Account, error) {
	var account domain.Account
	r.logger.Debug("Getting account", zap.String("username", username))
	err := r.pool.QueryRow(context.Background(), `
		SELECT id, username, password, created_at, updated_at
		FROM accounts WHERE username = LOWER($1)
	`, username).Scan(&account.ID, &account.Username, &account.Password, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}
