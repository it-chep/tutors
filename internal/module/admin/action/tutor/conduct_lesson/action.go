package conduct_lesson

import (
	"context"
	"fmt"
	"strconv"

	"github.com/it-chep/tutors.git/internal/module/admin/dto"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/conduct_lesson/dal"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// Action провести обычное занятие
type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, tutorID, studentID int64, durationInMinutes int64) error {
	// Получаем репетитора
	tutor, err := a.dal.GetTutor(ctx, tutorID)
	if err != nil {
		return err
	}

	// Получаем кошелек студента
	wallet, err := a.dal.GetStudentWallet(ctx, studentID)
	if err != nil {
		return err
	}

	// Вычисляем обновленное значение кошелька
	remain, err := a.getRemainBalance(tutor, wallet, durationInMinutes)
	if err != nil {
		return err
	}

	// ---- todo
	if remain.LessThan(decimal.NewFromInt(0)) {
		student, err := a.dal.GetStudent(ctx, studentID)
		if err != nil {
			return err
		}
		_ = student
		// пушим в бота сообщение студенту
	}
	// ---- todo

	// Помечаем урок проведенным
	err = a.dal.ConductLesson(ctx, tutorID, studentID, durationInMinutes)
	if err != nil {
		return err
	}

	return a.dal.UpdateStudentWallet(ctx, studentID, remain)
}

func (a *Action) getRemainBalance(tutor dto.Tutor, userWallet dto.Wallet, durationInMinutes int64) (decimal.Decimal, error) {
	costPerHour, err := strconv.ParseFloat(tutor.CostPerHour, 64)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid cost per hour: %w", err)
	}

	lessonCost := costPerHour * float64(durationInMinutes) / 60.0

	remain := userWallet.Balance.Sub(decimal.NewFromFloat(lessonCost))

	return remain, nil
}
