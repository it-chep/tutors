package check_auth

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/auth/check_auth/dal"
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

func (a *Action) Do(ctx context.Context, tutorID int64) error {
	return nil
}
