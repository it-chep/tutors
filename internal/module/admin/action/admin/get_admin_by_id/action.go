package get_admin_by_id

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/get_admin_by_id/dal"
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

func (a *Action) Do(ctx context.Context, adminID int64) error {
	return a.dal.GetAdminByID(ctx, adminID)
}
