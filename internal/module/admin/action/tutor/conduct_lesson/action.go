package conduct_lesson

import (
	"context"
	"fmt"
	"strconv"
	"time"

	accrualdal "github.com/it-chep/tutors.git/internal/module/admin/action/accrual/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/conduct_lesson/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
	"github.com/it-chep/tutors.git/internal/pkg/transaction"
	"github.com/it-chep/tutors.git/internal/pkg/transaction/wrapper"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// Action провести обычное занятие
type Action struct {
	dal        *dal.Repository
	accrualDal *accrualdal.Repository
	bot        *tg_bot.Bot
}

func New(pool *pgxpool.Pool, bot *tg_bot.Bot) *Action {
	db := wrapper.NewDatabase(pool)
	return &Action{
		dal:        dal.NewRepository(db),
		accrualDal: accrualdal.NewRepository(db),
		bot:        bot,
	}
}

func (a *Action) Do(ctx context.Context, tutorID, studentID, durationInMinutes int64, createdTime time.Time) error {
	// получаем студента
	student, err := a.dal.GetStudent(ctx, studentID)
	if err != nil {
		return err
	}

	// Получаем кошелек студента
	wallet, err := a.dal.GetStudentWallet(ctx, studentID)
	if err != nil {
		return err
	}

	// Вычисляем обновленное значение кошелька
	remain, err := a.getRemainBalance(student, wallet, durationInMinutes)
	if err != nil {
		return err
	}

	tutor, err := a.dal.GetTutor(ctx, tutorID)
	if err != nil {
		return err
	}

	accrualAmount, err := a.getTutorAccrual(tutor, durationInMinutes)
	if err != nil {
		return err
	}

	err = transaction.Exec(ctx, func(ctx context.Context) error {
		lessonID, err := a.dal.ConductLesson(ctx, tutorID, studentID, durationInMinutes, createdTime)
		if err != nil {
			return err
		}

		if err = a.dal.FinishTrial(ctx, studentID); err != nil {
			return err
		}

		if err = a.dal.UpdateStudentWallet(ctx, studentID, remain); err != nil {
			return err
		}

		return a.accrualDal.UpsertLessonAccrual(ctx, lessonID, tutorID, accrualAmount, createdTime)
	})
	if err != nil {
		return err
	}

	if remain.LessThan(decimal.Zero) {
		_ = a.bot.SendMessages([]bot_dto.Message{
			{
				Chat: student.ParentTgID,
				Text: "Добрый день! У вас возникла задолженность по оплате занятий, пополните пожалуйста баланс)",
			},
		})
	}

	return nil
}

func (a *Action) getRemainBalance(student dto.Student, userWallet dto.Wallet, durationInMinutes int64) (decimal.Decimal, error) {
	costPerHour, err := strconv.ParseFloat(student.CostPerHour, 64)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid cost per hour: %w", err)
	}

	lessonCost := costPerHour * float64(durationInMinutes) / 60.0

	remain := userWallet.Balance.Sub(decimal.NewFromFloat(lessonCost))

	return remain, nil
}

func (a *Action) getTutorAccrual(tutor dto.Tutor, durationInMinutes int64) (decimal.Decimal, error) {
	costPerHour, err := decimal.NewFromString(tutor.CostPerHour)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid tutor cost per hour: %w", err)
	}

	return costPerHour.Mul(decimal.NewFromInt(durationInMinutes)).Div(decimal.NewFromInt(60)), nil
}
