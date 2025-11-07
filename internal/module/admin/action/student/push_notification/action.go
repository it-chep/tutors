package push_notification

import (
	"context"
	"fmt"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/push_notification/dal"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

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

func (a *Action) Do(ctx context.Context, studentID int64) error {
	student, err := a.dal.GetStudentByID(ctx, studentID)
	if err != nil {
		return err
	}

	wallet, err := a.dal.GetStudentWallet(ctx, studentID)
	if err != nil {
		return err
	}

	if !wallet.Balance.LessThan(decimal.Zero) || wallet.Balance.Equal(decimal.Zero) {
		return errors.New("У пользователя нет задолженности")
	}

	err = a.bot.SendMessages([]bot_dto.Message{
		{
			Chat: student.ParentTgID,
			Text: fmt.Sprintf("Здравствуйте, у вас возникла задолженность по занятиям - %s.р, пополните пожалуйста баланс.", wallet.Balance.Abs().String()),
		},
	})
	if err != nil {
		return err
	}

	err = a.dal.AddNotificationToHistory(ctx, userCtx.UserIDFromContext(ctx), student.ParentTgID)
	if err != nil {
		return err
	}

	return nil
}
