package get_tutor_finance

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/get_tutor_finance/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"time"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, tutorID int64, from, to string) (dto.TutorFinance, error) {
	fromTime, err := time.Parse(time.DateTime, from)
	if err != nil {
		return dto.TutorFinance{}, errors.New("Неправильно указан формат даты 'ОТ'")
	}

	toTime, err := time.Parse(time.DateTime, to)
	if err != nil {
		return dto.TutorFinance{}, errors.New("Неправильно указан формат даты 'ДО'")
	}

	if toTime.Before(fromTime) {
		return dto.TutorFinance{}, errors.New("'ДО' раньше 'ОТ'")
	}

	if toTime.After(time.Now()) {
		return dto.TutorFinance{}, errors.New("'ДО' раньше чем сейчас")
	}

	return a.dal.GetFinanceInfo(ctx, tutorID, fromTime, toTime)
}
