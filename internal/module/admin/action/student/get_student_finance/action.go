package get_student_finance

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_student_finance/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
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

func (a *Action) Do(ctx context.Context, studentID int64, from, to string) (dto.StudentFinance, error) {
	fromTime, toTime, err := convert.StringsIntervalToTime(from, to)
	if err != nil {
		return dto.StudentFinance{}, err
	}

	return a.dal.GetFinanceInfo(ctx, studentID, fromTime, toTime)
}
