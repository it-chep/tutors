package tutor_by_id

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/tutor_by_id/dal"
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

func (a *Action) Do(ctx context.Context, tutorID int64) (tutor dto.Tutor, err error) {
	return a.dal.GetTutor(ctx, tutorID)
}
