package top_up_balance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/it-chep/tutors.git/internal/config"
	top_up_balance_dal "github.com/it-chep/tutors.git/internal/module/bot/action/commands/top_up_balance/dal"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/module/bot/dto/business"
	alfa "github.com/it-chep/tutors.git/internal/pkg/alpha"
	alfadto "github.com/it-chep/tutors.git/internal/pkg/alpha/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/it-chep/tutors.git/internal/pkg/tbank"
	tbankDto "github.com/it-chep/tutors.git/internal/pkg/tbank/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	alfa        *alfa.Client
	tbank       *tbank.Client
	bot         *tg_bot.Bot
	dal         *top_up_balance_dal.Dal
	bankByAdmin map[int64]config.Bank
}

func NewAction(pool *pgxpool.Pool, alfa *alfa.Client, tbank *tbank.Client, bankByAdmin map[int64]config.Bank, bot *tg_bot.Bot) *Action {
	return &Action{
		alfa:        alfa,
		tbank:       tbank,
		bot:         bot,
		dal:         top_up_balance_dal.NewDal(pool),
		bankByAdmin: bankByAdmin,
	}
}

func (a *Action) TransactionExists(ctx context.Context, msg dto.Message) bool {
	_, err := a.dal.TransactionByParent(ctx, msg.User)
	return err == nil
}

func (a *Action) InitTransaction(ctx context.Context, msg dto.Message) error {
	if !a.TransactionExists(ctx, msg) {
		_, err := a.dal.InitTransaction(ctx, msg.User)
		if err != nil {
			return err
		}
	}

	return a.bot.SendMessages([]bot_dto.Message{
		{Chat: msg.ChatID, Text: "Пожалуйста, укажите сумму на которую хотите пополнить баланс"},
	})
}

func (a *Action) SetAmount(ctx context.Context, msg dto.Message) error {
	amount, err := strconv.Atoi(msg.Text)
	if err != nil || amount <= 0 {
		return a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "Необходимо ввести целое положительное число"},
		})
	}

	transaction, err := a.dal.TransactionByParent(ctx, msg.User)
	if err != nil {
		return err
	}

	if err = a.dal.SetTransactionAmount(ctx, transaction.ID, int64(amount)); err != nil {
		logger.Error(ctx, "ошибка при установке суммы пополнения", err)
		return err
	}

	adminID, err := a.dal.AdminIDByParent(ctx, msg.User)
	if err != nil {
		logger.Error(ctx, "ошибка при получении админа от тутора родителя", err)
		return a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "Извините, но мы не нашли вашего репетитора. Обратитесь за помощью в поддержку"},
		})
	}

	var orderID, url string
	switch a.bankByAdmin[adminID] {
	case config.Alpha:
		orderID, url, err = a.regOrderAlpha(ctx, msg, adminID, transaction, amount)
	case config.TBank:
		orderID, url, err = a.regOrderTbank(ctx, msg, adminID, transaction, amount)
	}

	if orderID == "" || err != nil {
		return err
	}
	if err = a.dal.SetOrderID(ctx, transaction.ID, orderID); err != nil {
		logger.Error(ctx, "ошибка при сохранении идентификатора заказа", err)
		return err
	}

	return a.bot.SendMessages([]bot_dto.Message{
		{
			Chat: msg.ChatID,
			Text: fmt.Sprintf("Спасибо за использование нашего сервиса!\n Ждем вашей оплаты: %s \n\n⚠️ Рекомендуем оплачивать по СБП", url),
		},
	})
}

func (a *Action) regOrderAlpha(ctx context.Context, msg dto.Message, adminID int64, tx *business.Transaction, amount int) (orderID, url string, _ error) {
	resp, err := a.alfa.RegisterOrder(ctx, alfadto.NewOrderRequest(adminID, tx.ID, amount))
	if err != nil {
		if resp != nil {
			err = fmt.Errorf("%s: %s", err.Error(), resp.ErrorMessage)
		}
		logger.Error(ctx, "ошибка при создании платежки в альфабанке", err)
		if err = a.dal.DropTransaction(ctx, tx.ID); err != nil {
			logger.Error(ctx, "ошибка при удалении транзакции при ошибке от альфабанка", err)
			return "", "", err
		}
		return "", "", a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "У банка возникли технические неполадки, пожалуйста, попробуйте чуть позже"},
		})
	}
	return resp.OrderID, resp.FormURL, nil
}

func (a *Action) regOrderTbank(ctx context.Context, msg dto.Message, adminID int64, tx *business.Transaction, amount int) (orderID, url string, _ error) {
	resp, err := a.tbank.InitPayment(ctx, tbankDto.NewInitRequest(adminID, tx.ID, int64(amount)))
	if err != nil {
		if resp != nil {
			err = fmt.Errorf("%s: %s", err.Error(), resp.Message)
		}
		logger.Error(ctx, "ошибка при создании платежки в т банке", err)
		if err = a.dal.DropTransaction(ctx, tx.ID); err != nil {
			logger.Error(ctx, "ошибка при удалении транзакции при ошибке от т банка", err)
			return "", "", err
		}
		return "", "", a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "У банка возникли технические неполадки, пожалуйста, попробуйте чуть позже"},
		})
	}
	return resp.OrderID, resp.PaymentURL, nil
}
