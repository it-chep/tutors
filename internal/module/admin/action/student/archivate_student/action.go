package archivate_student

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/archivate_student/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, studentID int64) error {
	if dto.IsAdminRole(ctx) || dto.IsAssistantRole(ctx) {
		return a.adminUnArchivate(ctx, studentID)
	}

	return nil
}

func (a *Action) adminUnArchivate(ctx context.Context, studentID int64) error {
	reqAdminID := userCtx.AdminIDFromContext(ctx)

	adminID, err := a.dal.GetStudentAdminID(ctx, studentID)
	if err != nil {
		return err
	}

	if adminID != reqAdminID {
		return errors.New("invalid admin")
	}

	student, err := a.dal.GetStudent(ctx, studentID)
	if err != nil {
		return err
	}

	if student.IsArchived {
		return errors.New("student already archived")
	}

	err = a.dal.ArchivateStudent(ctx, studentID)
	if err != nil {
		return err
	}

	return nil
}
