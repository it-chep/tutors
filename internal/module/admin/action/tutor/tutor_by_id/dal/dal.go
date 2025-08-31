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

// GetTutor получение репетитора по ID
func (r *Repository) GetTutor(ctx context.Context, tutorID int64) (dto.Tutor, error) {
	sql := `
		select * from tutors where id = $1
	`

	args := []interface{}{
		tutorID,
	}

	var tutor dao.TutorDAO
	err := pgxscan.Get(ctx, r.pool, &tutor, sql, args)
	if err != nil {
		return dto.Tutor{}, err
	}

	return tutor.ToDomain(), nil
}
