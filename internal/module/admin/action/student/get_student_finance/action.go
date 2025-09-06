package get_student_finance

import (
	"context"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_student_finance/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
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
	fromTime, err := time.Parse(time.DateTime, from)
	if err != nil {
		return dto.StudentFinance{}, errors.New("Неправильно указан формат даты 'ОТ'")
	}

	toTime, err := time.Parse(time.DateTime, to)
	if err != nil {
		return dto.StudentFinance{}, errors.New("Неправильно указан формат даты 'ДО'")
	}

	if toTime.Before(fromTime) {
		return dto.StudentFinance{}, errors.New("'ДО' раньше 'ОТ'")
	}

	if toTime.After(time.Now()) {
		return dto.StudentFinance{}, errors.New("'ДО' раньше чем сейчас")
	}

	return a.dal.GetFinanceInfo(ctx, studentID, fromTime, toTime)
}
