package get_all_finance

import (
	"context"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Action struct {
	dal *dal.Repository
}

// TODO TODO TODO
func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, from, to string) (dto.StudentFinance, error) {
	fromTime, err := time.Parse(time.DateOnly, from)
	if err != nil {
		return dto.StudentFinance{}, errors.New("Неправильно указан формат даты 'ОТ'")
	}

	toTime, err := time.Parse(time.DateOnly, to)
	if err != nil {
		return dto.StudentFinance{}, errors.New("Неправильно указан формат даты 'ДО'")
	}
	// Используем текущую временную зону
	loc := time.Now().Location()

	// Устанавливаем время для fromTime: 00:00:00 в текущей локации
	fromTime = time.Date(
		fromTime.Year(),
		fromTime.Month(),
		fromTime.Day(),
		0, 0, 0, 0,
		loc,
	)

	// Устанавливаем время для toTime: 23:59:59 в текущей локации
	toTime = time.Date(
		toTime.Year(),
		toTime.Month(),
		toTime.Day(),
		23, 59, 59, 0,
		loc,
	)
	if toTime.Before(fromTime) {
		return dto.StudentFinance{}, errors.New("'ДО' раньше 'ОТ'")
	}

	if toTime.After(time.Now()) {
		return dto.StudentFinance{}, errors.New("'ДО' раньше чем сейчас")
	}

	return a.dal.GetAllFinanceInfo(ctx, fromTime, toTime)
}
