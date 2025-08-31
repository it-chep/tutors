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

func (a *Action) Do(ctx context.Context, studentID int64) error {
	// todo надо ли хранить время конца триалки либо будем смотреть по первой оплате и дате создания
	return a.dal.MarkStudentTrialDone(ctx, studentID)
}
