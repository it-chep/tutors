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

// DeleteStudent удаление студента
func (r *Repository) DeleteStudent(ctx context.Context, studentID int64) error {
	sql := `
		delete from students where id = $1
	`
	result, err := r.pool.Exec(ctx, sql, studentID)
	if result.RowsAffected() == 0 {
		return errors.New("Ошибка при удалении студена")
	}

	return err
}

// DeleteWallet удаление кошелька студента
func (r *Repository) DeleteWallet(ctx context.Context, studentID int64) error {
	sql := `
		delete from wallet where student_id = $1
	`
	_, err := r.pool.Exec(ctx, sql, studentID)
	return err
}
