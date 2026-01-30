package bot

import (
	"github.com/it-chep/tutors.git/internal/config"
	dtoInternal "github.com/it-chep/tutors.git/internal/dto"
	"github.com/it-chep/tutors.git/internal/module/bot/action"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Bot struct {
	Actions *action.Agg
}

func New(pool *pgxpool.Pool, config *config.Config, bot *tg_bot.Bot, gateways *dtoInternal.PaymentGateways) *Bot {
	return &Bot{
		Actions: action.NewAgg(pool, bot, config, gateways),
	}
}
