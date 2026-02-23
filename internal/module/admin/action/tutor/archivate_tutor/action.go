package archivate_tutor

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/archivate_tutor/dal"
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

func (a *Action) Do(ctx context.Context, tutorID int64) error {
	if dto.IsAdminRole(ctx) || dto.IsAssistantRole(ctx) {
		return a.archivate(ctx, tutorID)
	}
	return nil
}

func (a *Action) archivate(ctx context.Context, tutorID int64) error {
	reqAdminID := userCtx.AdminIDFromContext(ctx)

	adminID, err := a.dal.GetTutorAdminID(ctx, tutorID)
	if err != nil {
		return err
	}

	if adminID != reqAdminID {
		return errors.New("invalid admin")
	}

	isArchived, err := a.dal.IsTutorArchived(ctx, tutorID)
	if err != nil {
		return err
	}

	if isArchived {
		return errors.New("tutor already archived")
	}

	return a.dal.ArchivateTutor(ctx, tutorID)
}
