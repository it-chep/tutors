package delete_tutor

import (
	"context"

	adminDal "github.com/it-chep/tutors.git/internal/module/admin/dal"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/delete_tutor/dal"
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

func (a *Action) Do(ctx context.Context, tutorID int64) error {
	tgAdminUsernameID, _ := a.dal.GetTutorTgID(ctx, tutorID)

	err := a.dal.DeleteTutor(ctx, tutorID)
	if err != nil {
		return err
	}

	err = a.dal.UpdateStudents(ctx, tutorID)
	if err != nil {
		return err
	}

	if tgAdminUsernameID != 0 {
		_ = a.commonDal.DeleteTgAdminUsername(ctx, tgAdminUsernameID)
	}

	return nil
}
