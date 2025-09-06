package search_student

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/search_student/dal"
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

func (a *Action) Do(ctx context.Context, query string) (_ []dto.Student, _ error) {
	students, err := a.dal.SearchStudent(ctx, query)
	if err != nil {
		return nil, err
	}

	if dto.IsTutorRole(ctx) {
		studentsForTutor := make([]dto.Student, 0, len(students))
		for _, student := range students {
			studentsForTutor = append(studentsForTutor, dto.Student{
				ID:         student.ID,
				FirstName:  student.FirstName,
				MiddleName: student.MiddleName,
			})
		}

		return studentsForTutor, nil
	}

	return students, nil
}
