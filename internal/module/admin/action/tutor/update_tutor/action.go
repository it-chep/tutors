package update_tutor

import (
	"context"

	adminDal "github.com/it-chep/tutors.git/internal/module/admin/dal"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/update_tutor/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/update_tutor/dto"
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

func (a *Action) Do(ctx context.Context, tutorID int64, upd dto.UpdateRequest) error {
	adminID, oldTgID, err := a.dal.GetTutorTgInfo(ctx, tutorID)
	if err != nil {
		return err
	}

	if upd.TgAdminUsername != "" {
		tgID, err := a.commonDal.AddTgAdminUsername(ctx, adminID, upd.TgAdminUsername)
		if err != nil {
			return err
		}
		upd.TgAdminUsernameID = tgID
	}

	if err := a.dal.UpdateTutor(ctx, tutorID, upd); err != nil {
		return err
	}

	if oldTgID != 0 && oldTgID != upd.TgAdminUsernameID {
		_ = a.commonDal.DeleteTgAdminUsername(ctx, oldTgID)
	}

	return nil
}
