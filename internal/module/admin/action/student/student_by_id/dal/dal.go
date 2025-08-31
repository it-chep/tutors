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
