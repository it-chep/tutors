package add_manual_transaction

import (
	"context"

	"github.com/google/uuid"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/add_manual_transaction/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
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

	return a.dal.AddManualTransaction(ctx, studentID, amount)
}
