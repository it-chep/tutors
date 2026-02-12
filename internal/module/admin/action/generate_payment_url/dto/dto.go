package dto

import "github.com/it-chep/tutors.git/internal/config"

type Agg struct {
	InternalTransactionUUID string
	Payment                 config.AdminPayment
	Amount                  int
	StudentID               int64
}
