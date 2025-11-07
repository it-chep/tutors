package get_tutor_finance

import (
	"context"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/get_tutor_finance/dal"
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

func (a *Action) Do(ctx context.Context, tutorID int64, from, to time.Time) (dto.TutorFinance, error) {
	amount, err := a.dal.GetFinanceInfo(ctx, tutorID, from, to)
	if err != nil {
		return dto.TutorFinance{}, err
	}

	tutorWages, tutorHours, err := a.dal.GetTutorFinanceInfo(ctx, tutorID, from, to)
	if err != nil {
		return dto.TutorFinance{}, err
	}

	return dto.TutorFinance{
		Amount:     amount,
		Wages:      tutorWages,
		HoursCount: tutorHours,
	}, nil
}
