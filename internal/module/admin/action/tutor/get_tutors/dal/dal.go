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

func (r *Repository) GetTutors(ctx context.Context) ([]dto.Tutor, error) {
	sql := `
		select * from tutors order by tutors.id asc
	`
	var tutors dao.TutorsDao
	if err := pgxscan.Select(ctx, r.pool, &tutors, sql); err != nil {
		return nil, err
	}

	return tutors.ToDomain(), nil
}
