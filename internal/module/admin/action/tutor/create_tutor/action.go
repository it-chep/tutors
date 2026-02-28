package create_tutor

import (
	"context"

	adminDal "github.com/it-chep/tutors.git/internal/module/admin/dal"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/create_tutor/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/create_tutor/dto"
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

func (a *Action) Do(ctx context.Context, createDTO dto.Request, adminID int64) error {
	if createDTO.TgAdminUsername != "" {
		tgID, err := a.commonDal.AddTgAdminUsername(ctx, adminID, createDTO.TgAdminUsername)
		if err != nil {
			return err
		}
		createDTO.TgAdminUsernameID = tgID
	}

	tutorID, err := a.dal.CreateTutor(ctx, createDTO, adminID)
	if err != nil {
		return err
	}

	return a.dal.CreateUser(ctx, createDTO, tutorID, adminID)
}
