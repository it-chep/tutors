package delete_student

import (
	"context"

	adminDal "github.com/it-chep/tutors.git/internal/module/admin/dal"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/delete_student/dal"
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

func (a *Action) Do(ctx context.Context, studentID int64) error {
	tgAdminUsernameID, _ := a.dal.GetStudentTgID(ctx, studentID)

	err := a.dal.DeleteWallet(ctx, studentID)
	if err != nil {
		return err
	}
	if err := a.dal.DeleteStudent(ctx, studentID); err != nil {
		return err
	}

	if tgAdminUsernameID != 0 {
		_ = a.commonDal.DeleteTgAdminUsername(ctx, tgAdminUsernameID)
	}

	return nil
}
