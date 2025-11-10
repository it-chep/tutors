package start

import (
	"context"

	start_dal "github.com/it-chep/tutors.git/internal/module/bot/action/commands/start/dal"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	GetBalance   = "Мой баланс"
	TopUpBalance = "Пополнить баланс"
	GetLessons   = "Мои занятия"
)

type Action struct {
	bot *tg_bot.Bot
	dal *start_dal.Dal
}

func NewAction(pool *pgxpool.Pool, bot *tg_bot.Bot) *Action {
	return &Action{
		bot: bot,
		dal: start_dal.NewDal(pool),
	}
}

func (a *Action) Start(ctx context.Context, msg dto.Message) error {
	known, err := a.dal.IsKnown(ctx, msg.User)
	if err != nil {
		return err
	}

	if !known {
		return a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "Привет, кажется мы с вами пока что не знакомы, давайте это исправим! Как вас зовут? Напишите свои Фамилию Имя Отчество"},
		})
	}

	return a.bot.SendMessages([]bot_dto.Message{
		{
			Chat: msg.ChatID,
			Text: "Привет, чем могу вам помочь сегодня?",
			Buttons: dto.StepButtons{
				{Text: GetBalance},
				{Text: TopUpBalance},
				{Text: GetLessons},
			},
		},
	})
}
