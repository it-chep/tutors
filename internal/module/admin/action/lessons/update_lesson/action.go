package update_lesson

import (
	"context"
	"strconv"
	"time"

	accrualdal "github.com/it-chep/tutors.git/internal/module/admin/action/accrual/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
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

func (a *Action) Do(ctx context.Context, lessonID int64, upd dto.UpdateLesson) error {
	lesson, err := a.dal.GetLessonByID(ctx, lessonID)
	if err != nil {
		return err
	}

	student, err := a.dal.GetStudentInfo(ctx, lesson.StudentID)
	if err != nil {
		return err
	}

	tutor, err := a.dal.GetTutorInfo(ctx, lesson.TutorID)
	if err != nil {
		return err
	}

	newAccrualAmount, err := a.calculateTutorLessonCost(tutor.CostPerHour, upd.Duration)
	if err != nil {
		return err
	}

	return transaction.Exec(ctx, func(ctx context.Context) error {
		if upd.Duration != lesson.Duration {
			wallet, err := a.dal.GetStudentWallet(ctx, lesson.StudentID)
			if err != nil {
				return err
			}

			balance := wallet.Balance
			oldLessonCost := a.calculateStudentLessonCost(student.CostPerHour, lesson.Duration)
			newLessonCost := a.calculateStudentLessonCost(student.CostPerHour, upd.Duration)
			newBalance := balance.Add(oldLessonCost).Sub(newLessonCost)

			if err = a.dal.UpdateStudentBalance(ctx, student.ID, newBalance); err != nil {
				return err
			}
		}

		if err = a.dal.UpdateLesson(ctx, lessonID, upd); err != nil {
			return err
		}

		return a.accrualDal.UpsertLessonAccrual(ctx, lessonID, lesson.TutorID, newAccrualAmount, upd.Date)
	})
}

func (a *Action) calculateStudentLessonCost(studentCostPerHour string, lessonDuration time.Duration) decimal.Decimal {
	intCost, err := strconv.Atoi(studentCostPerHour)
	if err != nil {
		return decimal.Decimal{}
	}

	costPerHour := decimal.NewFromInt(int64(intCost))
	duration := decimal.NewFromFloat(lessonDuration.Hours())
	totalCost := costPerHour.Mul(duration)

	return totalCost
}

func (a *Action) calculateTutorLessonCost(tutorCostPerHour string, lessonDuration time.Duration) (decimal.Decimal, error) {
	costPerHour, err := decimal.NewFromString(tutorCostPerHour)
	if err != nil {
		return decimal.Zero, err
	}

	duration := decimal.NewFromFloat(lessonDuration.Hours())
	return costPerHour.Mul(duration), nil
}
