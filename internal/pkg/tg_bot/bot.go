package tg_bot

//import (
//	"net/http"
//
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"github.com/it-chep/tutors.git/internal/module/bot/dto"
//	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
//	"github.com/samber/lo"
//)
//
//type Config interface {
//	WebhookURL() string
//	Token() string
//	UseWebhook() bool
//}
//
//type Bot struct {
//	bot        *tgbotapi.BotAPI
//	updates    tgbotapi.UpdatesChannel
//	useWebhook bool
//}
//
//func NewTgBot(cfg Config) (*Bot, error) {
//	bot, err := tgbotapi.NewBotAPI(cfg.Token())
//	if err != nil {
//		return nil, err
//	}
//
//	// Режим вебхуков
//	if cfg.UseWebhook() {
//		hook, _ := tgbotapi.NewWebhook(cfg.WebhookURL() + cfg.Token() + "/")
//		_, err = bot.Request(hook)
//		if err != nil {
//			return nil, err
//		}
//
//		_, err = bot.GetWebhookInfo()
//		if err != nil {
//			return nil, err
//		}
//
//		return &Bot{
//			bot:        bot,
//			useWebhook: true,
//		}, nil
//	}
//
//	// Режим поллинга
//	_, _ = bot.Request(tgbotapi.DeleteWebhookConfig{})
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//	updates := bot.GetUpdatesChan(u)
//
//	return &Bot{
//		bot:        bot,
//		updates:    updates,
//		useWebhook: false,
//	}, nil
//}
//
//func (b *Bot) HandleUpdate(r *http.Request) (*tgbotapi.Update, error) {
//	return b.bot.HandleUpdate(r)
//}
//
//func (b *Bot) GetUpdates() tgbotapi.UpdatesChannel {
//	return b.updates
//}
//
//func (b *Bot) GetUser(message dto.Message) (bot_dto.User, error) {
//	member, err := b.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
//		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
//			ChatID: message.ChatID,
//			UserID: message.User,
//		},
//	})
//	if err != nil {
//		return bot_dto.User{}, err
//	}
//
//	return bot_dto.User{
//		ID:       member.User.ID,
//		Name:     member.User.FirstName,
//		UserName: member.User.UserName,
//		IsAdmin:  member.IsCreator() || member.IsAdministrator(),
//	}, nil
//}
//
//func (b *Bot) SendMessage(msg bot_dto.Message, options ...MsgOption) error {
//	message := tgbotapi.NewMessage(msg.Chat, msg.Text)
//	if len(msg.Buttons) != 0 {
//		message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
//			lo.Map(msg.Buttons, func(b dto.StepButton, _ int) []tgbotapi.InlineKeyboardButton {
//				return tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(b.Text, b.Text))
//			})...,
//		)
//	}
//
//	for _, opt := range options {
//		message = opt(message)
//	}
//	_, err := b.bot.Send(message)
//	return err
//}
//
//func (b *Bot) SendMessages(messages []bot_dto.Message) error {
//	for _, msg := range messages {
//		if err := b.SendMessage(msg); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//// SendMessageWithContentType отправляет сообщение с media
//func (b *Bot) SendMessageWithContentType(msg bot_dto.Message) error {
//	var message tgbotapi.Chattable
//
//	if msg.ContentType == dto.Video {
//		videoMsg := tgbotapi.NewVideo(msg.Chat, tgbotapi.FileID(msg.MediaID))
//		videoMsg.Caption = msg.Text
//
//		message = videoMsg
//	}
//	if msg.ContentType == dto.Photo {
//		photoMsg := tgbotapi.NewPhoto(msg.Chat, tgbotapi.FileID(msg.MediaID))
//		photoMsg.Caption = msg.Text
//
//		message = photoMsg
//	}
//	if msg.ContentType == dto.Audio {
//		audioMsg := tgbotapi.NewAudio(msg.Chat, tgbotapi.FileID(msg.MediaID))
//		audioMsg.Caption = msg.Text
//
//		message = audioMsg
//	}
//	if msg.ContentType == dto.Document {
//		documentMsg := tgbotapi.NewDocument(msg.Chat, tgbotapi.FileID(msg.MediaID))
//		documentMsg.Caption = msg.Text
//
//		message = documentMsg
//	}
//	if msg.ContentType == dto.VideoNote {
//		videoNoteMsg := tgbotapi.NewVideoNote(msg.Chat, 0, tgbotapi.FileID(msg.MediaID))
//
//		message = videoNoteMsg
//	}
//	if msg.ContentType == dto.Voice {
//		voiceMsg := tgbotapi.NewVoice(msg.Chat, tgbotapi.FileID(msg.MediaID))
//
//		message = voiceMsg
//	}
//
//	// todo сделать кнопки ?
//
//	_, err := b.bot.Send(message)
//	return err
//}
