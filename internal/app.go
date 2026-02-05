package internal

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/it-chep/tutors.git/internal/config"
	dtoInternal "github.com/it-chep/tutors.git/internal/dto"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/bot"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/worker"
	"github.com/it-chep/tutors.git/internal/server"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
)

type Workers []worker.Worker

type App struct {
	config *config.Config
	pool   *pgxpool.Pool

	server          *server.Server
	bot             *tg_bot.Bot
	smtp            *smtp.ClientSmtp
	paymentGateways *dtoInternal.PaymentGateways

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
		initPaymentCredConf(ctx).
		initSmtp(ctx).
		initPaymentGateways(ctx).
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

	if !a.config.BotIsActive() || (a.config.UseWebhook() && a.config.BotIsActive()) {
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

				msg := dto.Message{
					User:   update.SentFrom().ID,
					Text:   getTxt(update),
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

func getTxt(update tgbotapi.Update) string {
	txt := ""

	if update.Message != nil {
		txt = update.Message.Text
		if txt == "/start" {
			return txt
		}
		if strings.Contains(txt, "/start ") {
			return txt[len("/start "):]
		}
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Data
	}

	return txt
}
