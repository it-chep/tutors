package get_comments

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_comments/dal"
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

func (a *Action) Do(ctx context.Context, studentID int64) ([]dto.Comment, error) {
	return a.dal.GetByStudentID(ctx, studentID)
}
