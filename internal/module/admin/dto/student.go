package dto

import "github.com/shopspring/decimal"

type Student struct {
	ID         int64
	FirstName  string
	LastName   string
	MiddleName string
	Phone      string
	Tg         string

	CostPerHour string
	SubjectID   int64
	TutorID     int64

	ParentFullName string
	ParentPhone    string
	ParentTg       string

	Balance decimal.Decimal

	HasButtons          bool
	IsFinishedTrial     bool
	IsOnlyTrialFinished bool
	IsBalanceNegative   bool
	IsNewbie            bool
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
