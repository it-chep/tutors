package add_available_tg

import (
	"context"

	adminDal "github.com/it-chep/tutors.git/internal/module/admin/dal"
	userCtx "github.com/it-chep/tutors.git/pkg/context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/assistant/add_available_tg/dal"
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

func (a *Action) Do(ctx context.Context, assistantID int64, tgAdminUsername string, existingID *int64) error {
	adminID := userCtx.AdminIDFromContext(ctx)

	var (
		tgID int64
		err  error
	)

	// Либо сетим существующий, либо создаем новый
	if existingID != nil {
		err = a.commonDal.ExistTgAdminUsernameID(ctx, adminID, *existingID)
		if err != nil {
			return err
		}
		tgID = *existingID
	} else {
		tgID, err = a.commonDal.AddTgAdminUsername(ctx, adminID, tgAdminUsername)
		if err != nil {
			return err
		}
	}

	return a.dal.AddAvailableTg(ctx, assistantID, tgID)
}
