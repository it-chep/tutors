package create_admin

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/create_admin/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/create_admin/dto"
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
	return a.dal.CreateAdmin(ctx, createDTO)
}
