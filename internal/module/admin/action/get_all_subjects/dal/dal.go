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

func (r *Repository) GetSubjects(ctx context.Context) ([]dto.Subject, error) {
	sql := `
		select * from subjects
	`

	var subjects dao.SubjectsDao
	if err := pgxscan.Select(ctx, r.pool, &subjects, sql); err != nil {
		return nil, err
	}

	return subjects.ToDomain(), nil
}
