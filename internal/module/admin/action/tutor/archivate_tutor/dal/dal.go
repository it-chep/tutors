package dal

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
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

// GetTutorAdminID получение admin_id репетитора
func (r *Repository) GetTutorAdminID(ctx context.Context, tutorID int64) (int64, error) {
	sql := `select admin_id from tutors where id = $1`
	var adminID int64
	err := pgxscan.Get(ctx, r.pool, &adminID, sql, tutorID)
	return adminID, err
}

// IsTutorArchived проверяем архивирован ли репетитор
func (r *Repository) IsTutorArchived(ctx context.Context, tutorID int64) (bool, error) {
	sql := `select coalesce(is_archive, false) from tutors where id = $1`
	var isArchive bool
	err := pgxscan.Get(ctx, r.pool, &isArchive, sql, tutorID)
	return isArchive, err
}

// ArchivateTutor архивирование репетитора
func (r *Repository) ArchivateTutor(ctx context.Context, tutorID int64) error {
	sql := `update tutors set is_archive = true where id = $1`
	_, err := r.pool.Exec(ctx, sql, tutorID)
	return err
}
