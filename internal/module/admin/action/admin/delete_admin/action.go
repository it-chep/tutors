package delete_admin

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/delete_admin/dal"
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
	_ = a.dal.DeleteWallets(ctx, adminID)
	_ = a.dal.DeleteStudents(ctx, adminID)
	_ = a.dal.DeleteTutors(ctx, adminID)
	err := a.dal.DeleteAdmin(ctx, adminID)
	if err != nil {
		return err
	}
	return nil
}
