package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

func createCsrfTokenTable(conn *pgx.Conn) error {
	fmt.Println("Creating csrf_tokens table")
	_, err := conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS csrf_tokens (
            id SERIAL PRIMARY KEY,
            token VARCHAR(64) NOT NULL,
            account_id INTEGER NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
            expires_at TIMESTAMP NOT NULL,
            FOREIGN KEY (account_id) REFERENCES accounts(id),
			UNIQUE (token, account_id)
        );
        CREATE INDEX IF NOT EXISTS idx_csrf_account_id ON csrf_tokens(account_id);
        CREATE INDEX IF NOT EXISTS idx_csrf_token ON csrf_tokens(token);
    `)
	return err
}

func CreateCsrfToken(conn *pgx.Conn, token string, accountID int32, expiresAt time.Time) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO csrf_tokens (token, account_id, expires_at)
		VALUES ($1, $2, $3)
	`, token, accountID, expiresAt)
	return err
}

func DeleteCsrfToken(conn *pgx.Conn, token string) error {
	_, err := conn.Exec(context.Background(), `
		DELETE FROM csrf_tokens WHERE token = $1
	`, token)
	return err
}
