package order_checker

import (
	"context"
	"time"

	"github.com/it-chep/tutors.git/internal/config"
	dtoInternal "github.com/it-chep/tutors.git/internal/dto"
	job_dal "github.com/it-chep/tutors.git/internal/module/admin/job/order_checker/dal"
	"github.com/it-chep/tutors.git/internal/module/bot/dto/business"
	alfaDto "github.com/it-chep/tutors.git/internal/pkg/alpha/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	tbankDto "github.com/it-chep/tutors.git/internal/pkg/tbank/dto"
	tochkaDto "github.com/it-chep/tutors.git/internal/pkg/tochka/dto"
	"github.com/samber/lo"
)

type TransactionChecker struct {
	dal            *job_dal.Repository
	gateways       *dtoInternal.PaymentGateways
	paymentByAdmin config.PaymentsByAdmin
}

func NewTransactionChecker(dal *job_dal.Repository, gateways *dtoInternal.PaymentGateways, paymentByAdmin config.PaymentsByAdmin) *TransactionChecker {
	return &TransactionChecker{dal: dal, gateways: gateways, paymentByAdmin: paymentByAdmin}
}
func (c *TransactionChecker) ConfirmOrder(ctx context.Context, orderID, terminal string) error {
	if !c.gateways.TBank.KnownTerminal(ctx, terminal) {
		return nil
	}

	return c.dal.UpdateBalance(ctx, orderID)
}

func (c *TransactionChecker) Start(ctx context.Context) {
	transactions, err := c.dal.GetUnconfirmedOrders(ctx)
	if err != nil {
		logger.Error(ctx, "ошибка получения непотдвержденных транзакций", err)
		return
	}

	if len(transactions) == 0 {
		return
	}

	studentIDs := lo.Map(transactions, func(item *business.Transaction, _ int) int64 {
		return item.StudentID
	})

	adminByStudent, err := c.dal.AdminIDByStudents(ctx, studentIDs)
	if err != nil {
		logger.Error(ctx, "ошибка получения админов по студентам в джобе", err)
		return
	}

	for _, transaction := range transactions {
		payment := c.paymentByAdmin.Payment(adminByStudent[transaction.StudentID], transaction.PaymentID)

		time.Sleep(300 * time.Millisecond)
		if payment.Bank == config.Alpha {
			status, err := c.gateways.Alfa.GetOrderStatus(ctx, alfaDto.NewStatusRequest(payment.PaymentID, lo.FromPtr(transaction.OrderID)))
			if err != nil {
				logger.Error(ctx, "ошибка получения статуса заказа", err)
			}
			if status.OrderStatus.Confirmed() {
				if err = c.dal.UpdateBalance(ctx, status.OrderNumber); err != nil {
					logger.Error(ctx, "ошибка обновления баланса пользователя", err)
					return
				}
			}
		}

		if payment.Bank == config.TBank {
			status, err := c.gateways.TBank.GetOrderStatus(ctx, tbankDto.NewGetOrderRequest(payment.PaymentID, lo.FromPtr(transaction.OrderID)))
			if err != nil {
				logger.Error(ctx, "ошибка получения статуса заказа", err)
			}
			if status.IsPaid() {
				if err = c.dal.UpdateBalance(ctx, status.OrderID); err != nil {
					logger.Error(ctx, "ошибка обновления баланса пользователя", err)
					return
				}
			}
			if status.Cancelled() {
				if err = c.dal.DropTransaction(ctx, lo.FromPtr(transaction.OrderID)); err != nil {
				}
			}
		}

		if payment.Bank == config.Tochka {
			status, err := c.gateways.Tochka.GetOrderStatus(ctx, tochkaDto.NewGetOrderRequest(payment.PaymentID, lo.FromPtr(transaction.OrderID)))
			if err != nil {
				logger.Error(ctx, "ошибка получения статуса заказа", err)
			}
			if status.IsPaid() {
				if err = c.dal.UpdateBalanceByOrderID(ctx, lo.FromPtr(transaction.OrderID)); err != nil {
					logger.Error(ctx, "ошибка обновления баланса пользователя", err)
					return
				}
			}
			if status.Expired() {
				if err = c.dal.DropTransaction(ctx, lo.FromPtr(transaction.OrderID)); err != nil {
				}
			}
		}
	}
	return
}

func (c *TransactionChecker) Stop() {}
