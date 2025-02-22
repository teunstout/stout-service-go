package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func createJishoTable(conn *pgx.Conn) error {
	fmt.Println("Creating jisho table")
	_, err := conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS jisho (
            id SERIAL PRIMARY KEY,
        );
    `)
	return err
}
