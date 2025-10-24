package bot

import (
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/bot/action"
	alfa "github.com/it-chep/tutors.git/internal/pkg/alpha"
	"github.com/it-chep/tutors.git/internal/pkg/tbank"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Bot struct {
	Actions *action.Agg
}

func New(pool *pgxpool.Pool, config *config.Config, bot *tg_bot.Bot, alfa *alfa.Client, tbank *tbank.Client) *Bot {
	return &Bot{
		Actions: action.NewAgg(pool, bot, config, alfa, tbank),
	}
}
