package get_all_finance

import (
	"context"
	"sync"

	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance/dto"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, from, to string, adminID int64) (dto.GetAllFinanceDto, error) {
	fromTime, toTime, err := convert.StringsIntervalToTime(from, to)
	if err != nil {
		return dto.GetAllFinanceDto{}, err
	}

	var (
		cashFlow     decimal.Decimal
		finance      decimal.Decimal
		lessonsCount indto.TutorLessons
		conversion   float64
		wg           = sync.WaitGroup{}
	)

	// Получаем общий оборот
	wg.Add(1)
	go func() {
		defer wg.Done()
		gCashFlow, gErr := a.dal.GetCashFlow(ctx, fromTime, toTime, adminID)
		if gErr != nil {
			logger.Error(ctx, "Ошибка при получении оборота", gErr)
			return
		}
		cashFlow = gCashFlow
	}()

	// Получаем расходы на зарплаты
	wg.Add(1)
	go func() {
		defer wg.Done()
		gfinance, gErr := a.dal.GetFinanceInfo(ctx, fromTime, toTime, adminID)
		if gErr != nil {
			logger.Error(ctx, "Ошибка при расходов на зп", gErr)
			return
		}
		finance = gfinance
	}()

	// Получаем оплаченные уроки
	wg.Add(1)
	go func() {
		defer wg.Done()
		gLessonsCount, gErr := a.dal.GetLessons(ctx, fromTime, toTime, adminID)
		if gErr != nil {
			logger.Error(ctx, "Ошибка при получении оплаченных уроков", gErr)
			return
		}
		lessonsCount = gLessonsCount
	}()

	// Получаем конверсию
	wg.Add(1)
	go func() {
		defer wg.Done()
		gConversion, gErr := a.dal.GetTutorsConversion(ctx, fromTime, toTime, adminID)
		if gErr != nil {
			logger.Error(ctx, "Ошибка при получении конверсии", gErr)
			return
		}
		conversion = gConversion
	}()

	wg.Wait()

	return dto.GetAllFinanceDto{
		Profit:            finance.String(),
		CashFlow:          cashFlow.String(),
		Conversion:        conversion,
		CountLessons:      lessonsCount.LessonsCount,
		CountBaseLessons:  lessonsCount.BaseCount,
		CountTrialLessons: lessonsCount.TrialCount,
	}, nil
}
