package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"stout.dev/jisho/internal/domain"
)

type JishoRepositoryInterface struct {
	conn   *pgx.Conn
	logger domain.Logger
}

func NewJishoRepository(connString string, l domain.Logger) *JishoRepositoryInterface {
	c, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		l.Error("Unable to connect to database", map[string]interface{}{"error": err})
		return nil
	}

	return &JishoRepositoryInterface{conn: c, logger: l}
}

func (r *JishoRepositoryInterface) SaveSearchHistory(mid int32, keyword string) error {
	r.logger.Info("Saving search history", map[string]interface{}{"mid": mid, "keyword": keyword})

	_, err := r.conn.Exec(context.Background(), `
		INSERT INTO member_search_history(member_id, search_term)
		VALUES ($1, $2)
	`, mid, keyword)

	if err != nil {
		r.logger.Error("Failed to save search history", map[string]interface{}{"mid": mid, "keyword": keyword, "error": err})
	}

	return nil
}
