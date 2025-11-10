package get_notification_history

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_notification_history/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, studentID int64, from, to time.Time) ([]dto.NotificationHistory, error) {
	return a.dal.GetNotificationsByRange(ctx, studentID, from, to)
}
