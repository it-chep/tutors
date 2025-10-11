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

func (r *Repository) MoveOneStudent(ctx context.Context, newTutorID, studentID int64) error {
	sql := `
		update students set tutor_id = $2 where id = $1
	`

	_, err := r.pool.Exec(ctx, sql, studentID, newTutorID)
	return err
}

func (r *Repository) MoveAllStudents(ctx context.Context, oldTutorID, newTutorID int64) error {
	sql := `
		update students set tutor_id = $2 where tutor_id = $1
	`

	_, err := r.pool.Exec(ctx, sql, oldTutorID, newTutorID)
	return err
}
