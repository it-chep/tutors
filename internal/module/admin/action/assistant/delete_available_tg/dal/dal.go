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

// DeleteAvailableTg удаляет пользователю тгшку по ID
func (r *Repository) DeleteAvailableTg(ctx context.Context, assistantID int64, tgAdminUsernameID int64) error {
	sql := `
		update assistant_tgs
			set available_tg_ids = array_remove(available_tg_ids, $2)
		where user_id = $1
	`
	args := []interface{}{
		assistantID,
		tgAdminUsernameID,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
