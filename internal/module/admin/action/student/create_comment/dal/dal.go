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

func (r *Repository) CreateComment(ctx context.Context, userID, studentID int64, text string) error {
	sql := `
		insert into comments (user_id, text, student_id)
		values ($1, $2, $3)
	`

	_, err := r.pool.Exec(ctx, sql, userID, text, studentID)
	return err
}
