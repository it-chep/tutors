package internal

import (
	"context"
	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/server"
	"github.com/it-chep/tutors.git/internal/server/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
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
		//Bot:   bot.New(a.pool, a.bot),
		Admin: admin.New(a.pool),
	}
	return a
}

func (a *App) initServer(context.Context) *App {
	// todo: в NewHandler передаем сервис для админки или бота
	h := handler.NewHandler()
	srv := server.New(h)
	a.server = srv
	return a
}
