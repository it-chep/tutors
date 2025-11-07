package get_tg_admins_usernames

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_tg_admins_usernames/dal"
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

func (a *Action) Do(ctx context.Context, adminID int64) ([]string, error) {
	return a.dal.GetUsernames(ctx, adminID)
}
