package delete_comment

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/delete_comment/dal"
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

func (a *Action) Do(ctx context.Context, studentID, commentID int64) error {
	return a.dal.Delete(ctx, studentID, commentID)
}
