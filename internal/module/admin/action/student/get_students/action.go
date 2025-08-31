package get_students

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_students/dal"
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

func (a *Action) Do(ctx context.Context, tutorID int64) (students []dto.Student, err error) {
	// todo доделать красный желтый зеленый и белый

	if tutorID == 0 {
		return a.dal.GetAllStudents(ctx)
	}

	return a.dal.GetTutorStudents(ctx, tutorID)
}
