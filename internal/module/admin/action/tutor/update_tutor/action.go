package update_tutor

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/update_tutor/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/update_tutor/dto"
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

func (a *Action) Do(ctx context.Context, tutorID int64, upd dto.UpdateRequest) error {
	return a.dal.UpdateTutor(ctx, tutorID, upd)
}
