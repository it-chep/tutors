package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
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

// GetStudent .
func (r *Repository) GetStudent(ctx context.Context, studentID int64) (dto.Student, error) {
	sql := `
		select * from students where id = $1
	`
	var student dao.StudentDAO
	err := pgxscan.Get(ctx, r.pool, &student, sql, studentID)
	if err != nil {
		return dto.Student{}, err
	}

	return student.ToDomain(), nil
}

// GetStudentAdminID .
func (r *Repository) GetStudentAdminID(ctx context.Context, studentID int64) (int64, error) {
	sql := `
		select t.admin_id from students s join tutors t on s.tutor_id = t.id where s.id = $1
	`
	var adminID int64
	err := pgxscan.Get(ctx, r.pool, &adminID, sql, studentID)
	if err != nil {
		return 0, err
	}

	return adminID, nil
}

// ArchivateStudent .
func (r *Repository) ArchivateStudent(ctx context.Context, studentID int64) error {
	sql := `update students set is_archive = true where id = $1`

	_, err := r.pool.Exec(ctx, sql, studentID)
	return err
}
