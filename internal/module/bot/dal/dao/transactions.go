package dao

import (
	"github.com/it-chep/tutors.git/internal/module/bot/dto/business"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/it-chep/tutors.git/pkg/xo"
	"github.com/samber/lo"
)

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
	return transaction
}
