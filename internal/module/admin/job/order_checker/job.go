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
	"golang.org/x/sync/errgroup"
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

	grouppedByAdminPayment := make(map[int64]map[int64][]*business.Transaction)
	for _, transaction := range transactions {
		adminID := adminByStudent[transaction.StudentID]

		payment := c.paymentByAdmin.Payment(adminID, transaction.PaymentID)

		if len(grouppedByAdminPayment[adminID]) == 0 {
			grouppedByAdminPayment[adminID] = make(map[int64][]*business.Transaction, 3)
		}

		grouppedByAdminPayment[adminID][payment.PaymentID] = append(grouppedByAdminPayment[adminID][payment.PaymentID], transaction)
	}

	g, gCtx := errgroup.WithContext(ctx)
	for adminID, adminPayments := range grouppedByAdminPayment {

		adminID := adminID
		adminPayments := adminPayments

		g.Go(func() (gErr error) {
			for paymentID, paymentTransactions := range adminPayments {
				var (
					payment = c.paymentByAdmin.Payment(adminID, paymentID)
				)

				switch payment.Bank {
				case config.Alpha:
					gErr = c.ProcessAlfa(gCtx, payment, paymentTransactions)
				case config.TBank:
					gErr = c.ProcessTbank(gCtx, payment, paymentTransactions)
				case config.Tochka:
					gErr = c.ProcessTochka(gCtx, payment, paymentTransactions)
				}

				if gErr != nil {
					return err
				}
			}
			return gErr
		})

	}

	if err = g.Wait(); err != nil {
		logger.Error(ctx, "обработка транзакций завершена ошибкой", err)
	}
}

func (c *TransactionChecker) ProcessAlfa(ctx context.Context, payment config.AdminPayment, transactions []*business.Transaction) error {
	for _, transaction := range transactions {
		status, err := c.gateways.Alfa.GetOrderStatus(ctx, alfaDto.NewStatusRequest(payment.PaymentID, lo.FromPtr(transaction.OrderID)))
		if err != nil {
			if status.ErrorCode == "6" {
				_ = c.dal.DropTransaction(ctx, transaction.ID)
				continue
			}
		}
		if status.OrderStatus.Confirmed() {
			if err = c.dal.UpdateBalance(ctx, status.OrderNumber); err != nil {
				logger.Error(ctx, "ошибка обновления баланса пользователя", err)
				return err
			}
		}
		if status.OrderStatus.Cancelled() {
			_ = c.dal.DropTransaction(ctx, transaction.ID)
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (c *TransactionChecker) ProcessTbank(ctx context.Context, payment config.AdminPayment, transactions []*business.Transaction) error {
	for _, transaction := range transactions {
		status, err := c.gateways.TBank.GetOrderStatus(ctx, tbankDto.NewGetOrderRequest(payment.PaymentID, lo.FromPtr(transaction.OrderID)))
		if err != nil {
			logger.Error(ctx, "ошибка получения статуса заказа", err)
		}
		if status.IsPaid() {
			if err = c.dal.UpdateBalance(ctx, status.OrderID); err != nil {
				logger.Error(ctx, "ошибка обновления баланса пользователя", err)
				return err
			}
		}
		if status.Cancelled() {
			_ = c.dal.DropTransaction(ctx, transaction.ID)
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (c *TransactionChecker) ProcessTochka(ctx context.Context, payment config.AdminPayment, transactions []*business.Transaction) error {
	for _, transaction := range transactions {
		status, err := c.gateways.Tochka.GetOrderStatus(ctx, tochkaDto.NewGetOrderRequest(payment.PaymentID, lo.FromPtr(transaction.OrderID)))
		if err != nil {
			logger.Error(ctx, "ошибка получения статуса заказа", err)
		}
		if status.IsPaid() {
			if err = c.dal.UpdateBalanceByOrderID(ctx, lo.FromPtr(transaction.OrderID)); err != nil {
				logger.Error(ctx, "ошибка обновления баланса пользователя", err)
				return err
			}
		}
		if status.Expired() {
			_ = c.dal.DropTransaction(ctx, transaction.ID)
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (c *TransactionChecker) Stop() {}
