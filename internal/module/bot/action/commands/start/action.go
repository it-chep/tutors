package start

import (
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
)

const (
	GetBalance   = "Мой баланс"
	TopUpBalance = "Пополнить баланс"
)

type Action struct {
	bot *tg_bot.Bot
}

func NewAction(bot *tg_bot.Bot) *Action {
	return &Action{
		bot: bot,
	}
}

func (a *Action) Start(msg dto.Message) error {
	return a.bot.SendMessages([]bot_dto.Message{
		{
			Chat: msg.ChatID,
			Text: "Привет, чем могу вам помочь сегодня?",
			Buttons: dto.StepButtons{
				{Text: GetBalance},
				{Text: TopUpBalance},
			},
		},
	})
}
