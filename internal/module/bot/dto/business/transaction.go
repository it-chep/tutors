package business

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	// OrderNumber
	ID string
	// OrderID
	OrderID     *string
	CreatedAt   time.Time
	ConfirmedAt *time.Time
	Amount      *decimal.Decimal
	StudentID   int64
}
