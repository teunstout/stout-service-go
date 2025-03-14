package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func CreateLogin(conn *pgx.Conn, username string, password string) error {
	// TODO: Hash password https://auth0.com/blog/hashing-passwords-one-way-road-to-security/
	_, err := conn.Exec(context.Background(), `
		INSERT INTO accounts (username, password)
		VALUES (LOWER($1), LOWER($2))
	`, username, password)
	if err != nil {
		return fmt.Errorf("unable to create login: %v", err)
	}
	return nil
}

func AuthenticateUser(conn *pgx.Conn, username string, password string) (string, error) {
	var id string
	if err := conn.QueryRow(context.Background(), `
        SELECT id FROM login WHERE LOWER(username) = LOWER($1) AND password = $2
    `, username, password).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("User not found")
		}
		return "", err
	}

	return id, nil
}
