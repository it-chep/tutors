package conduct_lesson

import (
	"context"
	"fmt"

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

func (a *Action) Do(ctx context.Context, tutorID, studentID int64, duration int64) error {
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

	// ---- todo
	// Из кошелька вычитаем стоимость занятия
	//remain := wallet.Balance - tutor.CostPerHour*duration
	_ = tutor.CostPerHour
	if wallet.Balance.LessThan(decimal.NewFromInt(0)) {
		student, err := a.dal.GetStudent(ctx, studentID)
		if err != nil {
			return err
		}
		_ = student
		// пушим в бота сообщение студенту
	}
	// ---- todo

	return a.dal.UpdateStudentWallet(ctx, studentID, fmt.Sprintf("%s", wallet.Balance.String()))
}
