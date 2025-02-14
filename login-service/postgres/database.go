package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

const DATABASE_URL = "user=golang password=golang host=127.0.0.1 port=5432 dbname=golang sslmode=disable"

func Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return conn, nil
}

func InitDatabase(conn *pgx.Conn) error {
	if err := createAccountsTable(conn); err != nil {
		return err
	}

	return nil
}
