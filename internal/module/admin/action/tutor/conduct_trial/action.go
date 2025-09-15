package conduct_trial

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/conduct_trial/dal"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Action провести пробное занятие
type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, tutorID, studentID int64) error {

	// запоминаем что триалка проведена
	err := a.dal.ConductTrialLesson(ctx, studentID, tutorID)
	if err != nil {
		return err
	}
	return a.dal.MarkStudentTrialDone(ctx, studentID)
}
