package internal

import (
	"context"
	"fmt"
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type App struct {
	config *config.Config
	pool   *pgxpool.Pool

	server *server.Server
	bot    *tg_bot.Bot

	modules Modules
}

type Modules struct {
	//Bot   *bot.Bot
	Admin *admin.Module
}

func New(ctx context.Context) *App {
	cfg := config.NewConfig()

	app := &App{
		config: cfg,
	}

	app.initDB(ctx).
		initTgBot(ctx).
		initModules(ctx).
		initServer(ctx)

	return app
}

func (a *App) Run(ctx context.Context) {
	fmt.Println("start server http://localhost:8080")
	ctx = logger.ContextWithLogger(ctx, logger.New())

	if !a.config.BotIsActive() {
		log.Fatal(a.server.ListenAndServe())
	} else {
		go func() {
			log.Fatal(a.server.ListenAndServe())
		}()
	}

	//if !a.config.UseWebhook() && a.config.BotIsActive() {
	//	fmt.Println("Режим поллинга")
	//	// Режим поллинга
	//	for update := range a.bot.GetUpdates() {
	//		go func() {
	//			if update.ChatMember != nil {
	//				if lo.Contains([]string{"left", "kicked"}, update.ChatMember.NewChatMember.Status) {
	//					return
	//				}
	//				usrID := update.ChatMember.NewChatMember.User.ID
	//				chat := update.ChatMember.Chat.ID
	//				_ = a.modules.Bot.Actions.InvitePatient.InvitePatient(ctx, usrID, chat)
	//			}
	//
	//			if update.FromChat() == nil || update.SentFrom() == nil {
	//				return
	//			}
	//			logger.Message(ctx, "Обработка ивента")
	//
	//			txt := ""
	//			mediaID := ""
	//			if update.Message != nil {
	//				txt = update.Message.Text
	//				// фото
	//				if update.Message.Photo != nil {
	//					// массив фото разбивает фотографию на 4 качества, берем самое плохое )
	//					mediaID = update.Message.Photo[0].FileID
	//				}
	//				// видео
	//				if update.Message.Video != nil {
	//					mediaID = update.Message.Video.FileID
	//				}
	//				// документ
	//				if update.Message.Document != nil {
	//					mediaID = update.Message.Document.FileID
	//				}
	//				// кружок
	//				if update.Message.VideoNote != nil {
	//					mediaID = update.Message.VideoNote.FileID
	//				}
	//				// голосовое сообщение
	//				if update.Message.Voice != nil {
	//					mediaID = update.Message.Voice.FileID
	//				}
	//				// аудио сообщение
	//				if update.Message.Audio != nil {
	//					mediaID = update.Message.Audio.FileID
	//				}
	//			} else if update.CallbackQuery != nil {
	//				txt = update.CallbackQuery.Data
	//			}
	//
	//			msg := dto.Message{
	//				User:    update.SentFrom().ID,
	//				Text:    txt,
	//				ChatID:  update.FromChat().ID,
	//				MediaID: mediaID,
	//			}
	//			err := a.modules.Bot.Route(ctx, msg)
	//
	//			if err != nil {
	//				logger.Error(ctx, "Ошибка при обработке ивента", err)
	//			}
	//		}()
	//	}
	//}
}
