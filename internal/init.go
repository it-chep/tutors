package internal

import (
	"context"
	"log"
	"time"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	dtoInternal "github.com/it-chep/tutors.git/internal/dto"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/bot"
	"github.com/it-chep/tutors.git/internal/module/bot/dal"
	alfa "github.com/it-chep/tutors.git/internal/pkg/alpha"
	"github.com/it-chep/tutors.git/internal/pkg/alpha/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tbank"
	tbankDto "github.com/it-chep/tutors.git/internal/pkg/tbank/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tochka"
	tochkaDto "github.com/it-chep/tutors.git/internal/pkg/tochka/dto"
	"github.com/it-chep/tutors.git/internal/pkg/worker"
	"github.com/it-chep/tutors.git/internal/server"
	"github.com/it-chep/tutors.git/internal/server/handler"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {
	// ignore db columns that doesn't exist at the destination
	dbscanAPI, err := pgxscan.NewDBScanAPI(dbscan.WithAllowUnknownColumns(true))
	if err != nil {
		panic(err)
	}

	api, err := pgxscan.NewAPI(dbscanAPI)
	if err != nil {
		panic(err)
	}

	pgxscan.DefaultAPI = api
}

func (a *App) initDB(ctx context.Context) *App {
	pool, err := pgxpool.New(ctx, a.config.PgConn())
	if err != nil {
		log.Fatalf("[FATAL] не удалось создать кластер базы данных: %s", err)
	}

	a.pool = pool
	return a
}

func (a *App) initPaymentCredConf(ctx context.Context) *App {
	a.config.EnrichPayment(ctx, dal.NewDal(a.pool))
	return a
}

func (a *App) initSmtp(_ context.Context) *App {
	a.smtp = smtp.NewClientSmtp(a.config.SMTPConfig.Address, a.config.SMTPConfig.PassKey)
	return a
}

func (a *App) initPaymentGateways(_ context.Context) *App {
	paymentGates := &dtoInternal.PaymentGateways{
		Alfa: alfa.NewClient(dto.Credentials{
			BaseURL:         a.config.PaymentConfig.AlphaConf.BaseUrl,
			CredByPaymentID: a.config.PaymentConfig.AlphaConf.CredByID,
		}),
		TBank: tbank.NewClient(tbankDto.Credentials{
			BaseURL:         a.config.PaymentConfig.TBankConf.BaseUrl,
			CredByPaymentID: a.config.PaymentConfig.TBankConf.CredByID,
		}),
		Tochka: tochka.NewClient(tochkaDto.Credentials{
			BaseURL:         a.config.PaymentConfig.TochkaConf.BaseUrl,
			CredByPaymentID: a.config.PaymentConfig.TochkaConf.CredByID,
		}),
	}
	a.paymentGateways = paymentGates

	return a
}

func (a *App) initTgBot(_ context.Context) *App {
	if !a.config.BotIsActive() {
		return a
	}

	tgBot, err := tg_bot.NewTgBot(a.config)
	if err != nil {
		log.Fatal(err)
	}
	a.bot = tgBot
	return a
}

func (a *App) initModules(ctx context.Context) *App {
	a.modules = Modules{
		Bot:   bot.New(a.pool, a.config, a.bot, a.paymentGateways),
		Admin: admin.New(a.pool, a.smtp, a.config, a.bot, a.paymentGateways),
	}

	a.workers = append(a.workers, worker.NewWorker(ctx, a.modules.Admin.Checker.Start, 1*time.Minute, 1))
	return a
}

func (a *App) initServer(_ context.Context) *App {
	h := handler.NewHandler(a.bot, a.modules.Bot, a.modules.Admin, a.config)
	srv := server.New(h)
	a.server = srv
	return a
}
