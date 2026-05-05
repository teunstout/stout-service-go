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
	var exists bool
	err := r.conn.QueryRow(context.Background(), `
		SELECT EXISTS(SELECT 1 FROM accounts WHERE username = LOWER($1))
	`, username).Scan(&exists)
	return exists, err
}

func (r *AccountRepositoryInterface) CreateAccount(username string, password string) error {
	_, err := r.conn.Exec(context.Background(), `
		INSERT INTO accounts (username, password)
		VALUES (LOWER($1), $2)
	`, username, password)
	return err
}

func (r *AccountRepositoryInterface) GetAccount(username string) (domain.Account, error) {
	var account domain.Account
	err := r.conn.QueryRow(context.Background(), `
		SELECT id, username, password, created_at, updated_at
		FROM accounts WHERE username = LOWER($1)
	`, username).Scan(&account.ID, &account.Username, &account.Password, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}
