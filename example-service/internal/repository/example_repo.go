package repository

import (
	"context"
	"example-service/internal/domain"

	"github.com/jackc/pgx/v4"
)

type ExampleRepository struct {
	conn *pgx.Conn
}

func NewExampleRepository(conn *pgx.Conn) *ExampleRepository {
	return &ExampleRepository{conn: conn}
}

func (r *ExampleRepository) SaveExample(ctx context.Context, example domain.Example) error {
	_, err := r.conn.Exec(ctx, "INSERT INTO examples (id, name) VALUES ($1, $2)", example.ID, example.Name)
	return err
}
