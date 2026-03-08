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

func (r *Repository) Delete(ctx context.Context, studentID, commentID int64) error {
	sql := `
		delete from comments
		where id = $1 and student_id = $2
	`

	_, err := r.pool.Exec(ctx, sql, commentID, studentID)
	return err
}
