package dal

import (
	"context"

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

// GetTutorTgID получение tg_admin_username_id репетитора
func (r *Repository) GetTutorTgID(ctx context.Context, tutorID int64) (int64, error) {
	sql := `select coalesce(tg_admin_username_id, 0) from tutors where id = $1`
	var tgID int64
	err := r.pool.QueryRow(ctx, sql, tutorID).Scan(&tgID)
	return tgID, err
}

// DeleteTutor удаление репетитора
func (r *Repository) DeleteTutor(ctx context.Context, tutorID int64) (err error) {
	sql := `
		delete from tutors where id = $1
	`
	result, err := r.pool.Exec(ctx, sql, tutorID)
	if result.RowsAffected() == 0 {
		return errors.New("Ошибка при удалении репетитора")
	}

	return err
}

func (r *Repository) UpdateStudents(ctx context.Context, tutorID int64) (err error) {
	sql := `
		update students set tutor_id = 0 where tutor_id = $1
	`
	_, err = r.pool.Exec(ctx, sql, tutorID)

	return err
}
