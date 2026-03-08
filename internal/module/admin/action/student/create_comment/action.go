package create_comment

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/create_comment/dal"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
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

func (a *Action) Do(ctx context.Context, studentID int64, text string) error {
	return a.dal.CreateComment(ctx, userCtx.UserIDFromContext(ctx), studentID, text)
}
