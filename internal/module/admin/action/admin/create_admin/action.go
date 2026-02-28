package create_admin

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/create_admin/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/create_admin/dto"
	adminDal "github.com/it-chep/tutors.git/internal/module/admin/dal"
	dto2 "github.com/it-chep/tutors.git/internal/module/admin/dto"
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

func (a *Action) Do(ctx context.Context, createDTO dto.CreateRequest) error {
	userID, err := a.dal.CreateAdmin(ctx, createDTO)
	if err != nil {
		return err
	}

	if createDTO.Role == dto2.AssistantRole {
		return a.dal.AddAvailableTGs(ctx, userID, createDTO)
	}

	return nil
}
