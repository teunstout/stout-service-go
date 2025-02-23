package database

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

func CreateTables(conn *pgx.Conn) error {
	if err := createAccountsTable(conn); err != nil {
		return err
	}

	return nil
}

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
	if err != nil {
		return err
	}
	fmt.Println("Created accounts table")
	return nil
}

func CreateAccount(conn *pgx.Conn, username, password string) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO accounts (username, password)
		VALUES ($1, $2)
	`, username, password)
	if err != nil {
		return fmt.Errorf("unable to create login: %v", err)
	}
	return nil
}
