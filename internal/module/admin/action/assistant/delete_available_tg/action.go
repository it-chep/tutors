package delete_available_tg

import (
	"context"

	adminDal "github.com/it-chep/tutors.git/internal/module/admin/dal"
	userCtx "github.com/it-chep/tutors.git/pkg/context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/assistant/delete_available_tg/dal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	dal       *dal.Repository
	commonDal *adminDal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal:       dal.NewRepository(pool),
		commonDal: adminDal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, assistantID int64, tgAdminUsernameID int64) error {
	adminID := userCtx.AdminIDFromContext(ctx)

	err := a.commonDal.ExistTgAdminUsernameID(ctx, adminID, tgAdminUsernameID)
	if err != nil {
		return err
	}

	return a.dal.DeleteAvailableTg(ctx, assistantID, tgAdminUsernameID)
}
