package get_tutors

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/get_tutors/dal"
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

func (a *Action) Do(ctx context.Context) (tutors []dto.Tutor, err error) {
	return a.dal.GetTutors(ctx)
}
