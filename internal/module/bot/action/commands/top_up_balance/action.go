package top_up_balance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/it-chep/tutors.git/internal/config"
	dtoInternal "github.com/it-chep/tutors.git/internal/dto"
	top_up_balance_dal "github.com/it-chep/tutors.git/internal/module/bot/action/commands/top_up_balance/dal"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/module/bot/dto/business"
	alfadto "github.com/it-chep/tutors.git/internal/pkg/alpha/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	tbankDto "github.com/it-chep/tutors.git/internal/pkg/tbank/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot/bot_dto"
	tochkaDto "github.com/it-chep/tutors.git/internal/pkg/tochka/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	gateways       *dtoInternal.PaymentGateways
	bot            *tg_bot.Bot
	dal            *top_up_balance_dal.Dal
	paymentByAdmin config.PaymentsByAdmin
}

func NewAction(pool *pgxpool.Pool, gateways *dtoInternal.PaymentGateways, paymentByAdmin config.PaymentsByAdmin, bot *tg_bot.Bot) *Action {
	return &Action{
		gateways:       gateways,
		bot:            bot,
		dal:            top_up_balance_dal.NewDal(pool),
		paymentByAdmin: paymentByAdmin,
	}
}

func (a *Action) TransactionExists(ctx context.Context, msg dto.Message) bool {
	trans, err := a.dal.TransactionByParent(ctx, msg.User)
	return err == nil && trans != nil
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
	if transaction == nil {
		return a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "Пожалуйста, укажите сумму на которую хотите пополнить баланс"},
		})
	}

	if err = a.dal.SetTransactionAmount(ctx, transaction.ID, int64(amount)); err != nil {
		logger.Error(ctx, "ошибка при установке суммы пополнения", err)
		return err
	}

	adminID, paymentID, err := a.dal.AdminIDByParent(ctx, msg.User)
	if err != nil {
		logger.Error(ctx, "ошибка при получении админа от тутора родителя", err)
		return a.bot.SendMessages([]bot_dto.Message{
			{
				Chat: msg.ChatID,
				Text: "Извините, но мы не нашли вашего репетитора. Обратитесь за помощью в поддержку",
			},
		})
	}

	payment := a.paymentByAdmin.Payment(adminID, paymentID)
	var orderID, url string
	switch payment.Bank {
	case config.Alpha:
		orderID, url, err = a.regOrderAlpha(ctx, msg, payment.PaymentID, transaction, amount)
	case config.TBank:
		orderID, url, err = a.regOrderTbank(ctx, msg, payment.PaymentID, transaction, amount)
	case config.Tochka:
		orderID, url, err = a.regOrderTochka(ctx, msg, payment.PaymentID, transaction, amount)
	}

	if orderID == "" || err != nil {
		return err
	}
	if err = a.dal.SetOrderID(ctx, transaction.ID, orderID, payment.PaymentID); err != nil {
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

func (a *Action) regOrderAlpha(ctx context.Context, msg dto.Message, paymentID int64, tx *business.Transaction, amount int) (orderID, url string, _ error) {
	resp, err := a.gateways.Alfa.RegisterOrder(ctx, alfadto.NewOrderRequest(paymentID, tx.ID, amount))
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

func (a *Action) regOrderTbank(ctx context.Context, msg dto.Message, paymentID int64, tx *business.Transaction, amount int) (orderID, url string, _ error) {
	phone, err := a.dal.PhoneByStudent(ctx, tx.StudentID)
	if err != nil {
		return "", "", err
	}

	orderID, url, err = a.gateways.TBank.InitPayment(ctx, tbankDto.NewInitRequest(paymentID, tx.ID, int64(amount), phone))
	if err != nil {
		logger.Error(ctx, "ошибка при создании платежки в т банке", err)
		if err = a.dal.DropTransaction(ctx, tx.ID); err != nil {
			logger.Error(ctx, "ошибка при удалении транзакции при ошибке от т банка", err)
			return "", "", err
		}
		return "", "", a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "У банка возникли технические неполадки, пожалуйста, попробуйте чуть позже"},
		})
	}
	return orderID, url, nil
}

func (a *Action) regOrderTochka(ctx context.Context, msg dto.Message, paymentID int64, tx *business.Transaction, amount int) (orderID, url string, _ error) {
	resp, err := a.gateways.Tochka.InitPayment(ctx, tochkaDto.NewInitRequest(paymentID, int64(amount)))
	if err != nil {
		logger.Error(ctx, "ошибка при создании платежки в точке", err)
		if err = a.dal.DropTransaction(ctx, tx.ID); err != nil {
			logger.Error(ctx, "ошибка при удалении транзакции при ошибке от точки", err)
			return "", "", err
		}
		return "", "", a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "У банка возникли технические неполадки, пожалуйста, попробуйте чуть позже"},
		})
	}
	return resp.OperationID, resp.PaymentLink, nil
}
