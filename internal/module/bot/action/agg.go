package action

import (
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/acquaintance"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/get_balance"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/start"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/top_up_balance"
	alfa "github.com/it-chep/tutors.git/internal/pkg/alpha"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Agg struct {
	Start        *start.Action
	Acquaintance *acquaintance.Action
	GetBalance   *get_balance.Action
	TopUpBalance *top_up_balance.Action
}

func NewAgg(pool *pgxpool.Pool, bot *tg_bot.Bot, alfa *alfa.Client) *Agg {
	return &Agg{
		Start:        start.NewAction(pool, bot),
		Acquaintance: acquaintance.NewAction(pool, bot),
		GetBalance:   get_balance.NewAction(pool, bot),
		TopUpBalance: top_up_balance.NewAction(pool, alfa, bot),
	}
}
