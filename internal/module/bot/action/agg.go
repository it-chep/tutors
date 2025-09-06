package action

import (
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/get_balance"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/start"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Agg struct {
	Start      *start.Action
	GetBalance *get_balance.Action
}

func NewAgg(pool *pgxpool.Pool, bot *tg_bot.Bot) *Agg {
	return &Agg{
		Start:      start.NewAction(bot),
		GetBalance: get_balance.NewAction(pool, bot),
	}
}
