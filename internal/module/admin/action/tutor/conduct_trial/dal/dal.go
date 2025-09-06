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

func (r *Repository) MarkStudentTrialDone(ctx context.Context, studentID int64) error {
	sql := `
		update students set is_finished_trial = true where id = $1
	`

	_, err := r.pool.Exec(ctx, sql, studentID)
	return err
}
