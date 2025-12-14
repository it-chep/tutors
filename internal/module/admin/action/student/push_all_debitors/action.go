package push_all_debitors

import (
	"context"
	"fmt"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/push_all_debitors/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
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

func (a *Action) Do(ctx context.Context) error {
	students, err := a.dal.GetDebitors(ctx)
	if err != nil {
		return err
	}

	wallets, err := a.dal.GetStudentsWallet(ctx)
	if err != nil {
		return err
	}

	studentsWalletMap := lo.SliceToMap(wallets, func(item dto.Wallet) (int64, dto.Wallet) {
		return item.StudentID, item
	})

	for _, student := range students {
		debtStr := studentsWalletMap[student.ID].Balance.Abs().String()

		err = a.bot.SendMessages([]bot_dto.Message{
			{
				Chat: student.ParentTgID,
				Text: fmt.Sprintf("Здравствуйте, у вас возникла задолженность по занятиям: %sр, пополните пожалуйста баланс.", debtStr),
			},
		})
		if err != nil {
			return err
		}

		err = a.dal.AddNotificationsToHistory(ctx, userCtx.UserIDFromContext(ctx), student.ParentTgID)
		if err != nil {
			return err
		}

	}

	return nil
}
