package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func createAccountsTable(conn *pgx.Conn) error {
	fmt.Println("Creating accounts table")
	_, err := conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
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

func GetAccountPassword(conn *pgx.Conn, username string) (string, error) {
	var password string
	err := conn.QueryRow(context.Background(), `
		SELECT password FROM accounts WHERE username = LOWER($1)
	`, username).Scan(&password)
	return password, err
}
