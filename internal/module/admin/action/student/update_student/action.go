package update_student

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/update_student/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/update_student/dto"
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

func (a *Action) Do(ctx context.Context, studentID int64, upd dto.UpdateRequest) error {
	return a.dal.UpdateStudent(ctx, studentID, upd)
}
