package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// GetAssistantUsernames возвращает доступные тгшки ассистенту
func (r *Repository) GetAssistantUsernames(ctx context.Context, assistantID int64) ([]string, error) {
	sql := `
        SELECT available_tgs 
        FROM assistant_tgs 
        WHERE user_id = $1
    `

	var availableTgs []string
	err := pgxscan.Get(ctx, r.pool, &availableTgs, sql, assistantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}

	return availableTgs, nil
}
