package delete_tutor

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/delete_tutor/dal"
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

func (a *Action) Do(ctx context.Context, tutorID int64) error {
	err := a.dal.DeleteTutor(ctx, tutorID)
	if err != nil {
		return err
	}

	err = a.dal.UpdateStudents(ctx, tutorID)
	return err
}
