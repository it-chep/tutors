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
	// todo надо ли чистить кошелек и тд
	result, err := r.pool.Exec(ctx, sql, studentID)
	if result.RowsAffected() == 0 {
		return errors.New("Ошибка при удалении студена")
	}

	return err
}
