package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type CsrfRepositoryInterface struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewCsrfRepository(pool *pgxpool.Pool, l *zap.Logger) *CsrfRepositoryInterface {
	return &CsrfRepositoryInterface{
		pool:   pool,
		logger: l,
	}
}

func (r *CsrfRepositoryInterface) CreateCsrfToken(token string, accountID int32, expiresAt time.Time) error {
	r.logger.Debug("Creating CSRF token", zap.String("token", token), zap.Int32("account_id", accountID))
	_, err := r.pool.Exec(context.Background(), `
		INSERT INTO csrf_tokens (token, account_id, expires_at)
		VALUES ($1, $2, $3)
	`, token, accountID, expiresAt)
	return err
}

func (r *CsrfRepositoryInterface) GetCsrfIdBySessionToken(token string) (int32, error) {
	r.logger.Debug("Getting account ID by CSRF token", zap.String("token", token))
	var accountID int32
	err := r.pool.QueryRow(context.Background(), `
		SELECT account_id FROM csrf_tokens WHERE token = $1
	`, token).Scan(&accountID)
	return accountID, err
}

func (r *CsrfRepositoryInterface) DeleteCsrfTokensByAccountId(accountID int32) error {
	r.logger.Debug("Deleting CSRF tokens for account", zap.Int32("account_id", accountID))
	_, err := r.pool.Exec(context.Background(), `
		DELETE FROM csrf_tokens WHERE account_id = $1
	`, accountID)
	return err
}

func (r *CsrfRepositoryInterface) DeleteCsrfToken(token string) error {
	r.logger.Debug("Deleting CSRF token", zap.String("token", token))
	_, err := r.pool.Exec(context.Background(), `
		DELETE FROM csrf_tokens WHERE token = $1
	`, token)
	return err
}
