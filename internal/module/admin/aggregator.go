package admin

import (
	"github.com/it-chep/tutors.git/internal/config"
	dtoInternal "github.com/it-chep/tutors.git/internal/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/action"
	"github.com/it-chep/tutors.git/internal/module/admin/action/generate_payment_url"
	"github.com/it-chep/tutors.git/internal/module/admin/alpha"
	adminDal "github.com/it-chep/tutors.git/internal/module/admin/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/job/order_checker"
	job_dal "github.com/it-chep/tutors.git/internal/module/admin/job/order_checker/dal"
	tbankCallback "github.com/it-chep/tutors.git/internal/module/admin/tbank"
	"github.com/it-chep/tutors.git/internal/pkg/storage"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Module модуль отвечающий за работу админки
type Module struct {
	Actions   *action.Aggregator
	CommonDal *adminDal.Repository

	AlphaHook     *alpha.WebHookAlpha
	TbankCallback *tbankCallback.CallbackTbank
	Checker       *order_checker.TransactionChecker

	GeneratePaymentURL *generate_payment_url.Action
}

func New(
	pool *pgxpool.Pool, smtp *smtp.ClientSmtp, config *config.Config,
	bot *tg_bot.Bot, gateways *dtoInternal.PaymentGateways, objectStorage storage.Storage,
) *Module {
	actions := action.NewAggregator(pool, smtp, config.JwtConfig, bot, objectStorage)
	checker := order_checker.NewTransactionChecker(job_dal.NewRepository(pool), gateways, config.PaymentConfig.PaymentsByAdmin)
	return &Module{
		Actions:   actions,
		CommonDal: adminDal.NewRepository(pool),

		// TODO: уточнить секрет, точно ли альфа может передавать статичный Bearer?
		AlphaHook:     alpha.NewWebHookAlpha(checker, ""),
		TbankCallback: tbankCallback.NewCallbackTbank(checker),
		Checker:       checker,

		GeneratePaymentURL: generate_payment_url.NewAction(pool, gateways, config.PaymentConfig.PaymentsByAdmin),
	}
}
