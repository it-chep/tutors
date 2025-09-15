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

// ConductTrialLesson проводим триалку
func (r *Repository) ConductTrialLesson(ctx context.Context, studentID, tutorID int64) error {
	sql := `
		insert into conducted_lessons (student_id, tutor_id, is_trial, duration_in_minutes) 
		values ($1, $2, true, 0)
	`

	_, err := r.pool.Exec(ctx, sql, studentID, tutorID)
	return err
}

// MarkStudentTrialDone ставим пользователю, что он проходил триалку
func (r *Repository) MarkStudentTrialDone(ctx context.Context, studentID int64) error {
	sql := `
		update students set is_finished_trial = true where id = $1
	`

	_, err := r.pool.Exec(ctx, sql, studentID)
	return err
}
