package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type SessionTokenRepositoryInterface struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

func NewSessionTokenRepository(connString string, l *zap.Logger) *SessionTokenRepositoryInterface {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		l.Error("Unable to connect to database", zap.Error(err))
		return nil
	}

	return &SessionTokenRepositoryInterface{
		conn:   conn,
		logger: l,
	}
}

func (r *SessionTokenRepositoryInterface) CreateSessionToken(token string, accountID int32, expiresAt time.Time) error {
	_, err := r.conn.Exec(context.Background(), `
		INSERT INTO session_tokens (token, account_id, expires_at)
		VALUES ($1, $2, $3)
	`, token, accountID, expiresAt)
	return err
}

func (r *SessionTokenRepositoryInterface) GetAccountIdBySessionToken(token string) (int32, error) {
	var accountID int32
	err := r.conn.QueryRow(context.Background(), `
		SELECT account_id FROM session_tokens WHERE token = $1
	`, token).Scan(&accountID)
	return accountID, err
}

func (r *SessionTokenRepositoryInterface) DeleteSessionTokensByAccountId(accountID int32) error {
	_, err := r.conn.Exec(context.Background(), `
		DELETE FROM session_tokens WHERE account_id = $1
	`, accountID)
	return err
}

func (r *SessionTokenRepositoryInterface) DeleteSessionToken(token string) error {
	_, err := r.conn.Exec(context.Background(), `
		DELETE FROM session_tokens WHERE token = $1
	`, token)
	return err
}
