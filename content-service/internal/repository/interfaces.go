package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"stout.dev/content/internal/domain"
)

type TranslationRepository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	GetLists(ctx context.Context, accountID int32) ([]domain.TranslationListOutput, error)
	DeleteList(ctx context.Context, accountID int32, listID int32) (bool, error)
	DeleteEntries(ctx context.Context, accountID int32, ids []int32) ([]int32, error)
}
