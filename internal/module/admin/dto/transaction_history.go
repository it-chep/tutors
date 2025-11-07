package dto

import (
	"github.com/google/uuid"
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
}
