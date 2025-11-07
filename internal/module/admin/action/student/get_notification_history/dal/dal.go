package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetNotificationsByRange(ctx context.Context, studentID int64, from, to time.Time) ([]dto.NotificationHistory, error) {
	sql := `
		select nh.* 
		from notification_history nh  
		    join students s 
		        on nh.parent_tg_id = s.parent_tg_id 
		where s.id = $1
		order by nh.created_at desc
	`

	var history dao.NotificationsHistoryDAO
	err := pgxscan.Select(ctx, r.pool, &history, sql, studentID)

	return history.ToDomain(), err
}
