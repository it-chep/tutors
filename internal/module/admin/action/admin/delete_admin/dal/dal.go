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

// DeleteAdmin удаление админа
func (r *Repository) DeleteAdmin(ctx context.Context, adminID int64) error {
	sql := `
		delete from users where id = $1
	`
	_, err := r.pool.Exec(ctx, sql, adminID)
	return err
}

// UpdateTutorsAdmin обновление админов у репетиторов
func (r *Repository) UpdateTutorsAdmin(ctx context.Context, adminID int64) error {
	sql := `
		update tutors set admin_id = null where admin_id = $1
	`

	_, err := r.pool.Exec(ctx, sql, adminID)
	return err
}
