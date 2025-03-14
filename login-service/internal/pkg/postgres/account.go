package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

type Account struct {
	ID        int32
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func AccountExists(conn *pgx.Conn, username string) (bool, error) {
	var exists bool
	err := conn.QueryRow(context.Background(), `
		SELECT EXISTS(SELECT 1 FROM accounts WHERE username = LOWER($1))
	`, username).Scan(&exists)
	return exists, err
}

func CreateAccount(conn *pgx.Conn, username string, password string) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO accounts (username, password)
		VALUES (LOWER($1), $2)
	`, username, password)
	return err
}

func GetAccount(conn *pgx.Conn, username string) (Account, error) {
	var account Account
	err := conn.QueryRow(context.Background(), `
		SELECT id, username, password, created_at, updated_at
		FROM accounts WHERE username = LOWER($1)
	`, username).Scan(&account.ID, &account.Username, &account.Password, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}
