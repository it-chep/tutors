package info

import (
	"context"
	"fmt"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
)

// Action .
type Action struct {
	bot *tg_bot.Bot
}

// New ..
func New(bot *tg_bot.Bot) *Action {
	return &Action{
		bot: bot,
	}
}

// Do ..
func (a *Action) Do(_ context.Context, msg dto.Message) error {
	return a.bot.SendMessage(bot_dto.Message{
		Chat: msg.ChatID,
		Text: fmt.Sprintf("Ваши данные: \n\nTG_ID: %d", msg.ChatID),
	})
}
