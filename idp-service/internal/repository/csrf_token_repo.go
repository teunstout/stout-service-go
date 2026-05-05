package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type CsrfRepositoryInterface struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

func NewCsrfRepository(connString string, l *zap.Logger) *CsrfRepositoryInterface {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		l.Error("Unable to connect to database", zap.Error(err))
		return nil
	}

	return &CsrfRepositoryInterface{
		conn:   conn,
		logger: l,
	}
}

func (r *CsrfRepositoryInterface) CreateCsrfToken(token string, accountID int32, expiresAt time.Time) error {
	r.logger.Debug("Creating CSRF token", zap.String("token", token), zap.Int32("account_id", accountID))
	_, err := r.conn.Exec(context.Background(), `
		INSERT INTO csrf_tokens (token, account_id, expires_at)
		VALUES ($1, $2, $3)
	`, token, accountID, expiresAt)
	return err
}

func (r *CsrfRepositoryInterface) GetCsrfIdBySessionToken(token string) (int32, error) {
	r.logger.Debug("Getting account ID by CSRF token", zap.String("token", token))
	var accountID int32
	err := r.conn.QueryRow(context.Background(), `
		SELECT account_id FROM csrf_tokens WHERE token = $1
	`, token).Scan(&accountID)
	return accountID, err
}

func (r *CsrfRepositoryInterface) DeleteCsrfTokensByAccountId(accountID int32) error {
	r.logger.Debug("Deleting CSRF tokens for account", zap.Int32("account_id", accountID))
	_, err := r.conn.Exec(context.Background(), `
		DELETE FROM csrf_tokens WHERE account_id = $1
	`, accountID)
	return err
}

func (r *CsrfRepositoryInterface) DeleteCsrfToken(token string) error {
	r.logger.Debug("Deleting CSRF token", zap.String("token", token))
	_, err := r.conn.Exec(context.Background(), `
		DELETE FROM csrf_tokens WHERE token = $1
	`, token)
	return err
}
