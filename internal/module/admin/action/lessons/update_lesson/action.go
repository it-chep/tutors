package update_lesson

import (
	"context"
	"github.com/shopspring/decimal"
	"strconv"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Action провести обычное занятие
type Action struct {
	dal *dal.Repository
	bot *tg_bot.Bot
}

func New(pool *pgxpool.Pool, bot *tg_bot.Bot) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
		bot: bot,
	}
}

func (a *Action) Do(ctx context.Context, lessonID int64, upd dto.UpdateLesson) error {
	lesson, err := a.dal.GetLessonByID(ctx, lessonID)
	if err != nil {
		return err
	}

	// Изменяем кошелек, если изменили продолжительность
	if upd.Duration != lesson.Duration {
		student, err := a.dal.GetStudentInfo(ctx, lesson.StudentID)
		if err != nil {
			return err
		}

		wallet, err := a.dal.GetStudentWallet(ctx, lesson.StudentID)
		if err != nil {
			return err
		}

		balance := wallet.Balance

		oldLessonCost := a.calculateLessonCost(student.CostPerHour, lesson.Duration)
		newLessonCost := a.calculateLessonCost(student.CostPerHour, upd.Duration)

		newBalance := balance.Add(oldLessonCost).Sub(newLessonCost)

		err = a.dal.UpdateStudentBalance(ctx, student.ID, newBalance)
		if err != nil {
			return err
		}
	}

	err = a.dal.UpdateLesson(ctx, lessonID, upd)
	if err != nil {
		return err
	}

	return nil
}

func (a *Action) calculateLessonCost(studentCostPerHour string, lessonDuration time.Duration) decimal.Decimal {
	intCost, err := strconv.Atoi(studentCostPerHour)
	if err != nil {
		return decimal.Decimal{}
	}

	costPerHour := decimal.NewFromInt(int64(intCost))
	duration := decimal.NewFromFloat(lessonDuration.Hours())
	totalCost := costPerHour.Mul(duration)

	return totalCost
}
