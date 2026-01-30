package dao

import (
	"github.com/it-chep/tutors.git/internal/module/bot/dto/business"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/it-chep/tutors.git/pkg/xo"
	"github.com/samber/lo"
)

type TransactionDAOs []*TransactionDAO

func (daos TransactionDAOs) ToDomain() []*business.Transaction {
	return lo.Map(daos, func(dao *TransactionDAO, _ int) *business.Transaction {
		return dao.ToDomain()
	})
}

type TransactionDAO struct {
	xo.TransactionsHistory
}

func (s *TransactionDAO) ToDomain() *business.Transaction {
	transaction := &business.Transaction{
		ID:        s.ID.String(),
		CreatedAt: s.CreatedAt,
		StudentID: s.StudentID,
	}
	if s.ConfirmedAt.Valid {
		transaction.ConfirmedAt = lo.ToPtr(s.ConfirmedAt.Time)
	}
	if s.Amount.Valid {
		transaction.Amount = lo.ToPtr(convert.NumericToDecimal(s.Amount))
	}
	if s.OrderID.Valid {
		transaction.OrderID = lo.ToPtr(s.OrderID.String)
	}
	if s.PaymentID.Valid {
		transaction.PaymentID = s.PaymentID.Int64
	}
	return transaction
}
