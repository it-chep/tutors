package create_tutor

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/create_tutor/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/create_tutor/dto"
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

func (a *Action) Do(ctx context.Context, createDTO dto.Request, adminID int64) error {
	return a.dal.CreateTutor(ctx, createDTO, adminID)
}
