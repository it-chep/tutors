package bot

import (
	"github.com/it-chep/tutors.git/internal/module/bot/action"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Bot struct {
	Actions *action.Agg
}

func New(pool *pgxpool.Pool, bot *tg_bot.Bot) *Bot {
	return &Bot{
		Actions: action.NewAgg(pool, bot),
	}
}
