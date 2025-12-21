package dal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// AddAvailableTg добавлять пользователю тгшку
func (r *Repository) AddAvailableTg(ctx context.Context, assistantID int64, tgAdminUsername string) error {
	sql := `
		update assistant_tgs
			set available_tgs = array_append(available_tgs, $2)
		where user_id = $1
	`
	args := []interface{}{
		assistantID,
		tgAdminUsername,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
