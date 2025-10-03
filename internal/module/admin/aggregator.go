package admin

import (
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin/action"
	"github.com/it-chep/tutors.git/internal/module/admin/alpha"
	"github.com/it-chep/tutors.git/internal/module/admin/alpha/order_checker"
	alpha_dal "github.com/it-chep/tutors.git/internal/module/admin/alpha/order_checker/dal"
	alfa "github.com/it-chep/tutors.git/internal/pkg/alpha"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Module модуль отвечающий за работу админки
type Module struct {
	Actions *action.Aggregator

	AlphaHook *alpha.WebHookAlpha
	Checker   *order_checker.TransactionChecker
}

func New(pool *pgxpool.Pool, smtp *smtp.ClientSmtp, config config.JwtConfig, bot *tg_bot.Bot, client *alfa.Client) *Module {
	actions := action.NewAggregator(pool, smtp, config, bot)
	checker := order_checker.NewTransactionChecker(alpha_dal.NewRepository(pool), client)
	return &Module{
		Actions: actions,

		// TODO: уточнить секрет, точно ли альфа может передавать статичный Bearer?
		AlphaHook: alpha.NewWebHookAlpha(checker, ""),
		Checker:   checker,
	}
}
