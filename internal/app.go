package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/bot"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	alfa "github.com/it-chep/tutors.git/internal/pkg/alpha"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/worker"
	"github.com/it-chep/tutors.git/internal/server"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Workers []worker.Worker

type App struct {
	config *config.Config
	pool   *pgxpool.Pool

	server *server.Server
	bot    *tg_bot.Bot
	smtp   *smtp.ClientSmtp
	alfa   *alfa.Client

	modules Modules
	workers Workers
}

type Modules struct {
	Bot   *bot.Bot
	Admin *admin.Module
}

func New(ctx context.Context) *App {
	cfg := config.NewConfig()

	app := &App{
		config: cfg,
	}

	app.initDB(ctx).
		initSmtp(ctx).
		initAlfa(ctx).
		initTgBot(ctx).
		initModules(ctx).
		initServer(ctx)

	return app
}

func (a *App) Run(ctx context.Context) {
	fmt.Println("start server http://localhost:8080")
	ctx = logger.ContextWithLogger(ctx, logger.New())
	for _, w := range a.workers {
		w.Start(ctx)
	}

	if !a.config.BotIsActive() {
		log.Fatal(a.server.ListenAndServe())
	} else {
		go func() {
			log.Fatal(a.server.ListenAndServe())
		}()
	}

	if !a.config.UseWebhook() && a.config.BotIsActive() {
		fmt.Println("Режим поллинга")
		// Режим поллинга
		for update := range a.bot.GetUpdates() {
			go func() {
				if update.SentFrom() == nil ||
					update.FromChat() == nil {
					return
				}

				txt := ""
				if update.Message != nil {
					txt = update.Message.Text
				} else if update.CallbackQuery != nil {
					txt = update.CallbackQuery.Data
				}
				msg := dto.Message{
					User:   update.SentFrom().ID,
					Text:   txt,
					ChatID: update.FromChat().ID,
				}
				err := a.modules.Bot.Route(ctx, msg)

				if err != nil {
					logger.Error(ctx, "Ошибка при обработке ивента", err)
				}
			}()
		}
	}
}
