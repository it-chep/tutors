package change_all_students_payment

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/change_all_students_payment/dal"
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

// Do меняет платежку у всех студентов администратора
func (a *Action) Do(ctx context.Context, adminID, newPaymentID int64) error {
	return a.dal.ChangeAllPayments(ctx, adminID, newPaymentID)
}
