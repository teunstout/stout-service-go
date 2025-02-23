package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

const DATABASE_URL = "user=golang password=golang host=http://raspberrypi.local port=5432 dbname=development sslmode=disable"

func Connect() (*pgx.Conn, error) {
	connString := os.Getenv("CONNECTION_STRING")
	if connString == "" {
		connString = DATABASE_URL
	}

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return conn, nil
}

func InitDatabase(conn *pgx.Conn) error {
	if err := createAccountsTable(conn); err != nil {
		return err
	}
	if err := createSessionTokenTable(conn); err != nil {
		return err
	}
	if err := createCsrfTokenTable(conn); err != nil {
		return err
	}

	return nil
}
