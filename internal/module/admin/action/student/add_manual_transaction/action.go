package add_manual_transaction

import (
	"context"

	"github.com/google/uuid"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/add_manual_transaction/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/transaction"
	"github.com/it-chep/tutors.git/internal/pkg/transaction/wrapper"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(wrapper.NewDatabase(pool)),
	}
}

// Do создаёт ручную транзакцию для студента
func (a *Action) Do(ctx context.Context, studentID, amount int64) (uuid.UUID, error) {
	if dto.IsTutorRole(ctx) {
		return uuid.Nil, errors.New("access denied")
	}

	reqAdminID := userCtx.AdminIDFromContext(ctx)
	studentAdminID, err := a.dal.GetStudentAdminID(ctx, studentID)
	if err != nil {
		return uuid.Nil, err
	}

	if studentAdminID != reqAdminID {
		return uuid.Nil, errors.New("access denied")
	}

	wallet, err := a.dal.GetStudentWallet(ctx, studentID)
	if err != nil {
		return uuid.Nil, err
	}

	updatedBalance := wallet.Balance.Add(decimal.NewFromInt(amount))

	var transUUID uuid.UUID
	err = transaction.Exec(ctx, func(ctx context.Context) error {
		transUUID, err = a.dal.AddManualTransaction(ctx, studentID, amount)
		if err != nil {
			return err
		}
		return a.dal.UpdateStudentWallet(ctx, studentID, updatedBalance)
	})
	if err != nil {
		return uuid.Nil, err
	}

	return transUUID, nil
}
