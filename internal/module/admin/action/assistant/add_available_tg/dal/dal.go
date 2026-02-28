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

// AddAvailableTg добавлять пользователю тгшку по ID
func (r *Repository) AddAvailableTg(ctx context.Context, assistantID int64, tgAdminUsernameID int64) error {
	sql := `
		update assistant_tgs
			set available_tg_ids = array(
				select distinct unnest(array_append(available_tg_ids, $2))
			)
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
