package get_lessons

import (
	"context"
	"fmt"
	get_lessons_dal "github.com/it-chep/tutors.git/internal/module/bot/action/commands/get_lessons/dal"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"strings"
)

type Action struct {
	dal *get_lessons_dal.Dal
	bot *tg_bot.Bot
}

func NewAction(pool *pgxpool.Pool, bot *tg_bot.Bot) *Action {
	return &Action{
		dal: get_lessons_dal.NewDal(pool),
		bot: bot,
	}
}

func (a *Action) GetLessons(ctx context.Context, msg dto.Message) error {
	lessons, err := a.dal.GetLessons(ctx, msg.ChatID)
	if err != nil {
		return err
	}

	studentCost, err := a.dal.GetStudentCostByParentTgID(ctx, msg.ChatID)
	if err != nil {
		return err
	}

	if studentCost.IsZero() {
		return a.bot.SendMessages([]bot_dto.Message{
			{
				Chat: msg.ChatID,
				Text: "У вас пока нет занятий",
			},
		})
	}

	builder := strings.Builder{}
	builder.WriteString("Ваши занятия за последние 30 дней\n\n")

	for i, lesson := range lessons {
		lessonCost := studentCost.Mul(decimal.NewFromFloat(lesson.Duration.Hours()))
		builder.WriteString(
			fmt.Sprintf("%d. %s - %d мин - %s₽\n",
				i+1,
				fmt.Sprintf("%d.%d.%d", lesson.Date.Day(), lesson.Date.Month(), lesson.Date.Year()),
				int64(lesson.Duration.Minutes()),
				lessonCost.String(),
			))
	}

	return a.bot.SendMessages([]bot_dto.Message{
		{
			Chat: msg.ChatID,
			Text: builder.String(),
		},
	})
}
