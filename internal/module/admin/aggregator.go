package admin

import (
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin/action"
	"github.com/it-chep/tutors.git/internal/module/admin/alpha"
	alpha_dal "github.com/it-chep/tutors.git/internal/module/admin/alpha/dal"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Module модуль отвечающий за работу админки
type Module struct {
	Actions *action.Aggregator

	AlphaHook *alpha.WebHookAlpha
}

func New(pool *pgxpool.Pool, smtp *smtp.ClientSmtp, config config.JwtConfig, bot *tg_bot.Bot) *Module {
	actions := action.NewAggregator(pool, smtp, config)

	return &Module{
		Actions: actions,

		AlphaHook: alpha.NewWebHookAlpha(alpha_dal.NewRepository(pool)),
	}
}
