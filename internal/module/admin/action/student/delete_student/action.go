package delete_student

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/delete_student/dal"
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

func (a *Action) Do(ctx context.Context, studentID int64) error {
	// todo надо ли чистить кошелек и тд
	return a.dal.DeleteStudent(ctx, studentID)
}
