package permissions

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/assistant/permissions/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Action struct {
	dal *dal.Repository
}

type UpdateRequest struct {
	CanViewContracts      bool
	CanPenalizeAssistants []int64
}

type Permissions struct {
	CanViewContracts      bool
	CanPenalizeAssistants []int64
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Update(ctx context.Context, assistantID int64, req UpdateRequest) error {
	if !dto.IsAdminRole(ctx) && !dto.IsSuperAdminRole(ctx) {
		return errors.New("access denied")
	}

	adminID := userCtx.AdminIDFromContext(ctx)
	targetAdminID, err := a.dal.GetAssistantAdminID(ctx, assistantID)
	if err != nil {
		return err
	}
	if targetAdminID != adminID && !dto.IsSuperAdminRole(ctx) {
		return errors.New("access denied")
	}

	return a.dal.Update(ctx, assistantID, req.CanViewContracts, req.CanPenalizeAssistants)
}

func (a *Action) Get(ctx context.Context, assistantID int64) (Permissions, error) {
	permissions, err := a.dal.Get(ctx, assistantID)
	if err != nil {
		return Permissions{}, err
	}

	return Permissions{
		CanViewContracts:      permissions.CanViewContracts,
		CanPenalizeAssistants: permissions.CanPenalizeAssistants,
	}, nil
}
