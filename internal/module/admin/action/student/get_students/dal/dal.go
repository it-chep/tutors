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

func (r *Repository) GetAllStudents(ctx context.Context) ([]dto.Student, error) {
	sql := `
		select * from students
	`
	var students dao.StudentsDAO
	err := pgxscan.Select(ctx, r.pool, &students, sql)
	if err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}

func (r *Repository) GetTutorStudents(ctx context.Context, tutorID int64) ([]dto.Student, error) {
	sql := `
		select * from students where tutor_id = $1
	`
	var students dao.StudentsDAO
	err := pgxscan.Select(ctx, r.pool, &students, sql, tutorID)
	if err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}
