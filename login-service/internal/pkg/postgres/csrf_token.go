package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

func CreateCsrfToken(conn *pgx.Conn, token string, accountID int32, expiresAt time.Time) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO csrf_tokens (token, account_id, expires_at)
		VALUES ($1, $2, $3)
	`, token, accountID, expiresAt)
	return err
}

func GetCsrfIdBySessionToken(conn *pgx.Conn, token string) (int32, error) {
	var accountID int32
	err := conn.QueryRow(context.Background(), `
		SELECT account_id FROM csrf_tokens WHERE token = $1
	`, token).Scan(&accountID)
	return accountID, err
}

func DeleteCsrfTokensByAccountId(conn *pgx.Conn, accountID int32) error {
	_, err := conn.Exec(context.Background(), `
		DELETE FROM csrf_tokens WHERE account_id = $1
	`, accountID)
	return err
}

func DeleteCsrfToken(conn *pgx.Conn, token string) error {
	_, err := conn.Exec(context.Background(), `
		DELETE FROM csrf_tokens WHERE token = $1
	`, token)
	return err
}
