package dto

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"time"
)

type TransactionHistory struct {
	ID          uuid.UUID
	OrderID     string
	CreatedAt   time.Time
	IsConfirmed bool
	Amount      decimal.Decimal

	StudentID   int64
	StudentName string
}

type Transactions []TransactionHistory

func (t Transactions) StudentIDs() []int64 {
	return lo.Map(t, func(item TransactionHistory, _ int) int64 {
		return item.StudentID
	})
}
