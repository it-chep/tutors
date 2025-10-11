package move_students

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/move_students/dal"
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

func (a *Action) Do(ctx context.Context, oldTutorID, newTutorID, studentID int64) error {
	if studentID == 0 {
		return a.moveAllStudents(ctx, oldTutorID, newTutorID)
	}

	return a.moveOneStudent(ctx, newTutorID, studentID)
}

func (a *Action) moveAllStudents(ctx context.Context, oldTutorID, newTutorID int64) error {
	return a.dal.MoveAllStudents(ctx, oldTutorID, newTutorID)
}

func (a *Action) moveOneStudent(ctx context.Context, newTutorID, studentID int64) error {
	return a.dal.MoveOneStudent(ctx, newTutorID, studentID)
}
