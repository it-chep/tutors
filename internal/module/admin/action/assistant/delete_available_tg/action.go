package delete_available_tg

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/assistant/delete_available_tg/dal"
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

func (a *Action) Do(ctx context.Context, assistantID int64, tgAdminUsername string) error {
	return a.dal.DeleteAvailableTg(ctx, assistantID, tgAdminUsername)
}
