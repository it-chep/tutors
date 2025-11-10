package dto

import (
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type Student struct {
	ID         int64
	FirstName  string
	LastName   string
	MiddleName string
	Phone      string
	Tg         string

	CostPerHour string
	SubjectID   int64
	SubjectName string
	TutorID     int64
	TutorName   string

	ParentFullName  string
	ParentPhone     string
	ParentTg        string
	TgAdminUsername string

	Balance decimal.Decimal

	IsFinishedTrial     bool
	IsOnlyTrialFinished bool
	IsBalanceNegative   bool
	IsNewbie            bool

	ParentTgID int64
}

type Students []Student

func (s Students) IDs() []int64 {
	return lo.Map(s, func(item Student, _ int) int64 {
		return item.ID
	})
}

type StudentFinance struct {
	Count  int64
	Amount decimal.Decimal
}

type Wallet struct {
	ID        int64
	StudentID int64
	Balance   decimal.Decimal
}

type StudentWithTransactions struct {
	StudentID         int64
	TutorID           int64
	IsFinishedTrial   bool
	TransactionsCount int64
	Balance           decimal.Decimal
}
