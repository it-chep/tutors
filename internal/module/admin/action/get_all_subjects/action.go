package get_all_subjects

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_subjects/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context) ([]dto.Subject, error) {
	return a.dal.GetSubjects(ctx)
}
