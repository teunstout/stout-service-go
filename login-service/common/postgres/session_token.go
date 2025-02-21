package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

func createSessionTokenTable(conn *pgx.Conn) error {
	fmt.Println("Creating session_tokens table")
	_, err := conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS session_tokens (
            id SERIAL PRIMARY KEY,
            token VARCHAR(64) NOT NULL,
            account_id INTEGER NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
            expires_at TIMESTAMP NOT NULL,
            FOREIGN KEY (account_id) REFERENCES accounts(id),
			UNIQUE (token, account_id)
        );
        CREATE INDEX IF NOT EXISTS idx_session_account_id ON session_tokens(account_id);
        CREATE INDEX IF NOT EXISTS idx_session_token ON session_tokens(token);
    `)
	return err
}

func CreateSessionToken(conn *pgx.Conn, token string, accountID int32, expiresAt time.Time) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO session_tokens (token, account_id, expires_at)
		VALUES ($1, $2, $3)
	`, token, accountID, expiresAt)
	return err
}

func GetAccountIdBySessionToken(conn *pgx.Conn, token string) (int32, error) {
	var accountID int32
	err := conn.QueryRow(context.Background(), `
		SELECT account_id FROM session_tokens WHERE token = $1
	`, token).Scan(&accountID)
	return accountID, err
}

func DeleteSessionTokensByAccountId(conn *pgx.Conn, accountID int32) error {
	_, err := conn.Exec(context.Background(), `
		DELETE FROM session_tokens WHERE account_id = $1
	`, accountID)
	return err
}

func DeleteSessionToken(conn *pgx.Conn, token string) error {
	_, err := conn.Exec(context.Background(), `
		DELETE FROM session_tokens WHERE token = $1
	`, token)
	return err
}
