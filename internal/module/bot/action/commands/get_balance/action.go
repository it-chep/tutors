package get_balance

import (
	"context"
	"fmt"

	get_balance_dal "github.com/it-chep/tutors.git/internal/module/bot/action/commands/get_balance/dal"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	dal *get_balance_dal.Dal
	bot *tg_bot.Bot
}

func NewAction(pool *pgxpool.Pool, bot *tg_bot.Bot) *Action {
	return &Action{
		dal: get_balance_dal.NewDal(pool),
		bot: bot,
	}
}

func (a *Action) GetBalance(ctx context.Context, msg dto.Message) error {
	balance, err := a.dal.GetBalance(ctx, msg.User)
	if err != nil {
		return err
	}
	return a.bot.SendMessages([]bot_dto.Message{
		{
			Chat: msg.ChatID,
			Text: fmt.Sprintf(
				"На вашем балансе сейчас %s рублей",
				balance.String(),
			),
		},
	})
}
