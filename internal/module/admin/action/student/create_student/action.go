package create_student

import (
	"context"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/create_student/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/create_student/dto"
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

func (a *Action) Do(ctx context.Context, createDTO dto.CreateRequest) error {
	studentID, err := a.dal.CreateStudent(ctx, createDTO)
	if err != nil {
		return err
	}
	err = a.dal.CreateWallet(ctx, studentID)
	if err != nil {
		return err
	}

	if indto.IsAssistantRole(ctx) && len(createDTO.TgAdminUsername) != 0 {
		return a.dal.AddTgToAssistant(ctx, userCtx.UserIDFromContext(ctx), createDTO.TgAdminUsername)
	}

	return nil
}
