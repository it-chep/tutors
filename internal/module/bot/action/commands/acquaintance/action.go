package acquaintance

import (
	"context"

	acquaintance_dal "github.com/it-chep/tutors.git/internal/module/bot/action/commands/acquaintance/dal"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/start"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	bot *tg_bot.Bot
	dal *acquaintance_dal.Dal
}

func NewAction(pool *pgxpool.Pool, bot *tg_bot.Bot) *Action {
	return &Action{
		bot: bot,
		dal: acquaintance_dal.NewDal(pool),
	}
}

func (a *Action) OnRegistration(ctx context.Context, msg dto.Message) (onReg bool, _ error) {
	onReg, err := a.dal.ParentOnRegistration(ctx, msg.User)
	if err != nil {
		return false, err
	}

	return onReg, nil
}

func (a *Action) MakeKnown(ctx context.Context, msg dto.Message) error {
	ok, err := a.dal.MakeParentKnown(ctx, msg.User, msg.Text)
	if err != nil {
		return err
	}

	if !ok {
		return a.bot.SendMessages([]bot_dto.Message{
			{
				Chat: msg.ChatID,
				Text: `Не смог вас распознать, нажмите на /start и попробуйте еще раз. Если проблема не решилась, напишите своему администратору`,
			},
		})
	}

	return a.bot.SendMessages([]bot_dto.Message{
		{
			Chat: msg.ChatID,
			Text: "Отлично, будем знакомы! Чем могу помочь вам сегодня?",
			Buttons: dto.StepButtons{
				{Text: start.GetBalance},
				{Text: start.TopUpBalance},
				{Text: start.GetLessons},
			}},
	})
}
