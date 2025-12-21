package get_all_finance_by_tgs

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance_by_tgs/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance_by_tgs/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"sync"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, req dto.Request) (dto.GetAllFinanceDto, error) {
	var (
		cashFlow   decimal.Decimal
		finance    decimal.Decimal
		debt       decimal.Decimal
		tutorsInfo dto.TutorsInfo
		wg         = sync.WaitGroup{}
	)

	lessonsDaos, err := a.dal.GetLessonsInfo(ctx, req)
	if err != nil {
		return dto.GetAllFinanceDto{}, err
	}

	// Получаем общий оборот
	wg.Add(1)
	go func() {
		defer wg.Done()
		gCashFlow, gErr := a.dal.GetCashFlow(ctx, req, lessonsDaos)
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
		gfinance, gErr := a.dal.GetFinanceInfo(ctx, req, lessonsDaos)
		if gErr != nil {
			logger.Error(ctx, "Ошибка при расходов на зп", gErr)
			return
		}
		finance = gfinance
	}()

	// Получаем дебиторскую задолженность
	wg.Add(1)
	go func() {
		defer wg.Done()
		gdebt, gErr := a.dal.GetDebt(ctx, req)
		if gErr != nil {
			logger.Error(ctx, "Ошибка при получении дебиторской задолженности", gErr)
			return
		}
		debt = gdebt
	}()

	// Получаем дебиторскую задолженность
	wg.Add(1)
	go func() {
		defer wg.Done()
		gtutorsInfo, gErr := a.dal.GetTutorsInfo(ctx, req, lessonsDaos)
		if gErr != nil {
			logger.Error(ctx, "Ошибка при получении дебиторской задолженности", gErr)
			return
		}
		tutorsInfo = gtutorsInfo
	}()

	wg.Wait()

	return dto.GetAllFinanceDto{
		Profit:     finance.String(),
		CashFlow:   cashFlow.String(),
		Debt:       debt.String(),
		TutorsInfo: tutorsInfo,
	}, nil
}
