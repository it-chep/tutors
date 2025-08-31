package student_by_id

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/student_by_id/dal"
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

func (a *Action) Do(ctx context.Context, studentID int64) (dto.Student, error) {
	student, err := a.dal.GetStudent(ctx, studentID)
	// todo доделать шильдики о том что нет оплат или новичок

	return student, err
}
