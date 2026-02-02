package filter_students

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/filter_students/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/filter_students/dto"
	commondal "github.com/it-chep/tutors.git/internal/module/admin/dal"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Action struct {
	dal       *dal.Repository
	commonDal *commondal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal:       dal.NewRepository(pool),
		commonDal: commondal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, filter dto.FilterRequest) ([]indto.Student, error) {

	adminID := userCtx.AdminIDFromContext(ctx)
	students, err := a.dal.FilterStudents(ctx, adminID, filter)
	if err != nil {
		return nil, err
	}

	studentIDs := students.IDs()

	studentsWalletMap, err := a.dal.GetStudentsWallets(ctx, studentIDs)
	if err != nil {
		return nil, err
	}

	payments, err := a.dal.HasStudentsPayments(ctx, studentIDs)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении информации об оплатах студентов репетитора", err)
	}

	paymentsInfo, err := a.commonDal.GetStudentsPayments(ctx, studentIDs)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении информации о платежках студентов", err)
	}

	for i, _ := range students {
		// Задолженности
		wallet, ok := studentsWalletMap[students[i].ID]
		students[i].Balance = wallet.Balance
		students[i].Payment = paymentsInfo[students[i].ID]

		if !ok {
			students[i].IsBalanceNegative = false
		}
		students[i].IsBalanceNegative = wallet.Balance.LessThan(decimal.NewFromFloat(0.0))

		// Оплаты/новичок
		hasPayments, ok := payments[students[i].ID]
		if ok {
			students[i].IsNewbie = false
			students[i].IsOnlyTrialFinished = false
			continue
		}
		students[i].IsNewbie = !hasPayments && !students[i].IsFinishedTrial
		students[i].IsOnlyTrialFinished = !hasPayments

		students[i].Balance = wallet.Balance
	}

	return students, nil
}
