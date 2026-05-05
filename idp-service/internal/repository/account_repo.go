package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
)

type AccountRepositoryInterface struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

func NewAccountRepository(connString string, l *zap.Logger) *AccountRepositoryInterface {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		l.Error("Unable to connect to database", zap.Error(err))
		return nil
	}

	return &AccountRepositoryInterface{
		conn:   conn,
		logger: l,
	}
}

func (r *AccountRepositoryInterface) AccountExists(username string) (bool, error) {
	r.logger.Debug("Checking if account exists", zap.String("username", username))
	var exists bool
	err := r.conn.QueryRow(context.Background(), `
		SELECT EXISTS(SELECT 1 FROM accounts WHERE username = LOWER($1))
	`, username).Scan(&exists)
	return exists, err
}

func (r *AccountRepositoryInterface) CreateAccount(username string, password string) error {
	r.logger.Debug("Creating account", zap.String("username", username))
	_, err := r.conn.Exec(context.Background(), `
		INSERT INTO accounts (username, password)
		VALUES (LOWER($1), $2)
	`, username, password)
	return err
}

func (r *AccountRepositoryInterface) GetAccount(username string) (domain.Account, error) {
	var account domain.Account
	r.logger.Debug("Getting account", zap.String("username", username))
	err := r.conn.QueryRow(context.Background(), `
		SELECT id, username, password, created_at, updated_at
		FROM accounts WHERE username = LOWER($1)
	`, username).Scan(&account.ID, &account.Username, &account.Password, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}
