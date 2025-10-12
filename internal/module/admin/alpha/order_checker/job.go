package order_checker

import (
	"context"
	"fmt"

	alpha_dal "github.com/it-chep/tutors.git/internal/module/admin/alpha/order_checker/dal"
	"github.com/it-chep/tutors.git/internal/module/bot/dto/business"
	alfa "github.com/it-chep/tutors.git/internal/pkg/alpha"
	"github.com/it-chep/tutors.git/internal/pkg/alpha/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type TransactionChecker struct {
	dal  *alpha_dal.Repository
	alfa *alfa.Client
}

func NewTransactionChecker(dal *alpha_dal.Repository, alfa *alfa.Client) *TransactionChecker {
	return &TransactionChecker{dal: dal, alfa: alfa}
}

func (c *TransactionChecker) UpdateTransactionsByAmount(ctx context.Context, amount decimal.Decimal) error {
	transactions, err := c.dal.GetOrdersByAmount(ctx, amount)
	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		return nil
	}

	transactions = lo.Filter(transactions, func(item *business.Transaction, _ int) bool {
		if item.OrderID == nil {
			logger.Message(ctx, fmt.Sprintf("у транзакции %s нет номера заказа", item.ID))
			return false
		}
		return true
	})

	studentIDs := lo.Map(transactions, func(item *business.Transaction, _ int) int64 {
		return item.StudentID
	})

	adminByStudent, err := c.dal.AdminIDByStudents(ctx, studentIDs)
	if err != nil {
		logger.Error(ctx, "ошибка получения админов по студентам в вебхуке", err)
		return err
	}

	for _, transaction := range transactions {
		status, errG := c.alfa.GetOrderStatus(ctx, dto.NewStatusRequest(adminByStudent[transaction.StudentID], lo.FromPtr(transaction.OrderID)))
		if errG != nil {
			logger.Error(ctx, "ошибка получения статуса заказа", err)
			return errG
		}
		if status.OrderStatus.Confirmed() {
			if err = c.dal.UpdateBalance(ctx, status.OrderNumber); err != nil {
				logger.Error(ctx, "ошибка обновления баланса пользователя", err)
				return err
			}
		}
	}
	return nil
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
		status, err := c.alfa.GetOrderStatus(ctx, dto.NewStatusRequest(adminByStudent[transaction.StudentID], lo.FromPtr(transaction.OrderID)))
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
	return
}

func (c *TransactionChecker) Stop() {}
