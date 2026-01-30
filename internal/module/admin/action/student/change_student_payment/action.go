package change_student_payment

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/change_student_payment/dal"
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

// Do меняем платежку у студента
func (a *Action) Do(ctx context.Context, adminID, studentID, newPaymentID int64) error {
	studentAdminID, err := a.dal.GetStudentAdminID(ctx, studentID)
	if err != nil {
		return err
	}

	if studentAdminID != adminID {
		return errors.New("Доступ запрещен")
	}

	return a.dal.ChangePayment(ctx, studentID, newPaymentID)
}
