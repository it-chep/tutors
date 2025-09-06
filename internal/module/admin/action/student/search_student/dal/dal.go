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

func (r *Repository) SearchStudent(ctx context.Context, query string) ([]dto.Student, error) {
	sql := `
		select * from students
        where 
            concat(first_name, ' ', last_name, ' ', middle_name) ilike '%' || $1 || '%' or
            concat(last_name, ' ', first_name, ' ', middle_name) ilike '%' || $1 || '%' or
            parent_full_name ilike '%' || $1 || '%'
        order by id  
	`
	var students dao.StudentsDAO
	err := pgxscan.Select(ctx, r.pool, &students, sql, query)
	if err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}
