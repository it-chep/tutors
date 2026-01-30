package admin

import (
	"github.com/it-chep/tutors.git/internal/config"
	dtoInternal "github.com/it-chep/tutors.git/internal/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/action"
	"github.com/it-chep/tutors.git/internal/module/admin/alpha"
	"github.com/it-chep/tutors.git/internal/module/admin/job/order_checker"
	job_dal "github.com/it-chep/tutors.git/internal/module/admin/job/order_checker/dal"
	tbankCallback "github.com/it-chep/tutors.git/internal/module/admin/tbank"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Module модуль отвечающий за работу админки
type Module struct {
	Actions *action.Aggregator

	AlphaHook     *alpha.WebHookAlpha
	TbankCallback *tbankCallback.CallbackTbank
	Checker       *order_checker.TransactionChecker
}

func New(
	pool *pgxpool.Pool, smtp *smtp.ClientSmtp, config *config.Config,
	bot *tg_bot.Bot, gateways *dtoInternal.PaymentGateways,
) *Module {
	actions := action.NewAggregator(pool, smtp, config.JwtConfig, bot)
	checker := order_checker.NewTransactionChecker(job_dal.NewRepository(pool), gateways, config.PaymentConfig.PaymentsByAdmin)
	return &Module{
		Actions: actions,

		// TODO: уточнить секрет, точно ли альфа может передавать статичный Bearer?
		AlphaHook:     alpha.NewWebHookAlpha(checker, ""),
		TbankCallback: tbankCallback.NewCallbackTbank(checker),
		Checker:       checker,
	}
}
