package delete_lesson

import (
	"context"
	"fmt"
	"strconv"

	accrualdal "github.com/it-chep/tutors.git/internal/module/admin/action/accrual/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/delete_lesson/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/transaction"
	"github.com/it-chep/tutors.git/internal/pkg/transaction/wrapper"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// Action провести обычное занятие
type Action struct {
	dal        *dal.Repository
	accrualDal *accrualdal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	db := wrapper.NewDatabase(pool)
	return &Action{
		dal:        dal.NewRepository(db),
		accrualDal: accrualdal.NewRepository(db),
	}
}

func (a *Action) Do(ctx context.Context, lessonID int64) error {
	// получаем урок
	lesson, err := a.dal.GetLesson(ctx, lessonID)
	if err != nil {
		return err
	}

	// получаем студента
	student, err := a.dal.GetStudent(ctx, lesson.StudentID)
	if err != nil {
		return err
	}

	// Получаем кошелек студента
	wallet, err := a.dal.GetStudentWallet(ctx, student.ID)
	if err != nil {
		return err
	}

	// Вычисляем обновленное значение кошелька
	remain, err := a.getRemainBalance(student, wallet, lesson.Duration.Minutes())
	if err != nil {
		return err
	}

	return transaction.Exec(ctx, func(ctx context.Context) error {
		if err = a.dal.UpdateStudentWallet(ctx, student.ID, remain); err != nil {
			return err
		}

		if err = a.dal.DeleteLesson(ctx, lessonID); err != nil {
			return err
		}

		return a.accrualDal.DeleteLessonAccrual(ctx, lessonID)
	})
}

func (a *Action) getRemainBalance(student dto.Student, userWallet dto.Wallet, durationInMinutes float64) (decimal.Decimal, error) {
	costPerHour, err := strconv.ParseFloat(student.CostPerHour, 64)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid cost per hour: %w", err)
	}

	lessonCost := costPerHour * durationInMinutes / 60.0

	// добавляем к кошельку списанную сумму, так как урок добавлен ошибочно
	remain := userWallet.Balance.Add(decimal.NewFromFloat(lessonCost))

	return remain, nil
}
