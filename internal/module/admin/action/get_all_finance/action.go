package get_all_finance

import (
	"context"
	"sync"

	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance/dto"
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
		cashFlow decimal.Decimal
		finance  decimal.Decimal
		wg       = sync.WaitGroup{}
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

	wg.Wait()

	return dto.GetAllFinanceDto{
		Profit:   finance.String(),
		CashFlow: cashFlow.String(),
	}, nil
}
