package filter_students

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/filter_students/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/filter_students/dto"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
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

func (a *Action) Do(ctx context.Context, filter dto.FilterRequest) ([]indto.Student, error) {

	adminID := userCtx.UserIDFromContext(ctx)

	students, err := a.dal.FilterStudents(ctx, adminID, filter)
	if err != nil {
		return nil, err
	}

	studentsWalletMap, err := a.dal.GetStudentsWallets(ctx, students.IDs())
	if err != nil {
		return nil, err
	}

	for _, student := range students {
		wallet, ok := studentsWalletMap[student.ID]
		if !ok {
			continue
		}
		student.Balance = wallet.Balance
	}

	return students, nil
}
