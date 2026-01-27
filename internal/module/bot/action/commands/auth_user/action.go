package auth_user

import (
	"context"
	"fmt"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/auth_user/dal"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/start"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	dal *dal.Dal
	bot *tg_bot.Bot
}

func NewAction(pool *pgxpool.Pool, bot *tg_bot.Bot) *Action {
	return &Action{
		dal: dal.NewDal(pool),
		bot: bot,
	}
}

func (a *Action) Do(ctx context.Context, msg dto.Message, studentID int64) error {
	exist, err := a.dal.IsStudentExist(ctx, studentID)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("ошибка при получении студента: TGID: %d, StudID: %d ", msg.ChatID, studentID), err)
		return err
	}
	if !exist {
		return a.bot.SendMessages([]bot_dto.Message{
			{
				Chat: msg.ChatID,
				Text: fmt.Sprintf("Не могу найти для вас студента, напишите своему администратору. Ваш номер: %d", msg.ChatID),
			},
		})
	}

	attached, err := a.dal.IsParentAlreadyAttached(ctx, msg.ChatID)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("ошибка при получении родителя: TGID: %d, StudID: %d ", msg.ChatID, studentID), err)
		return err
	}
	if attached {
		return a.bot.SendMessages([]bot_dto.Message{
			{
				Chat: msg.ChatID,
				Text: fmt.Sprintf("Вы уже прикреплены к студенту, если хотите поменять студента, напишите своему администратору. Ваш номер: %d", msg.ChatID),
			},
		})
	}

	alreadyWithTG, err := a.dal.IsStudentAlreadyWithTgID(ctx, studentID)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("ошибка при получении информации о студенте: TGID: %d, StudID: %d ", msg.ChatID, studentID), err)
		return err
	}
	if alreadyWithTG {
		return a.bot.SendMessages([]bot_dto.Message{
			{
				Chat: msg.ChatID,
				Text: fmt.Sprintf("Студент недоступен к прикреплению, напишите своему администратору. Ваш номер: %d", msg.ChatID),
			},
		})
	}

	err = a.dal.AttachParentToStudent(ctx, studentID, msg.ChatID)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("ошибка при присоединении студента: TGID: %d, StudID: %d ", msg.ChatID, studentID), err)
		return a.bot.SendMessages([]bot_dto.Message{
			{
				Chat: msg.ChatID,
				Text: fmt.Sprintf("Не смог вас распознать, напишите своему администратору. Ваш номер: %d", msg.ChatID),
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
