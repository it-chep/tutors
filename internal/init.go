package internal

import (
	"context"
	"log"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/bot"
	alfa "github.com/it-chep/tutors.git/internal/pkg/alpha"
	"github.com/it-chep/tutors.git/internal/pkg/alpha/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
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

func (a *App) initSmtp(context.Context) *App {
	a.smtp = smtp.NewClientSmtp(a.config.SMTPConfig.Address, a.config.SMTPConfig.PassKey)
	return a
}

func (a *App) initAlfa(context.Context) *App {
	a.alfa = alfa.NewClient(dto.Credentials{
		BaseURL:  a.config.PaymentConfig.BaseUrl,
		UserName: a.config.PaymentConfig.User,
		Password: a.config.PaymentConfig.Password,
	})
	return a
}

func (a *App) initTgBot(context.Context) *App {
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

func (a *App) initModules(context.Context) *App {
	a.modules = Modules{
		Bot:   bot.New(a.pool, a.bot, a.alfa),
		Admin: admin.New(a.pool, a.smtp, a.config.JwtConfig),
	}
	return a
}

func (a *App) initServer(context.Context) *App {
	// todo: в NewHandler передаем сервис для админки или бота
	h := handler.NewHandler(a.modules.Admin)
	srv := server.New(h)
	a.server = srv
	return a
}
