package get_assistance

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/assistant/get_assistance/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
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

func (a *Action) Do(ctx context.Context, adminID int64) ([]dto.User, error) {
	return a.dal.GetAssistants(ctx, adminID)
}
