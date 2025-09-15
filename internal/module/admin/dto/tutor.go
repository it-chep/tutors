package dto

import "github.com/shopspring/decimal"

type Tutor struct {
	ID          int64
	FullName    string
	Phone       string
	Tg          string
	CostPerHour string
	SubjectID   int64
	SubjectName string
	AdminID     int64

	HasBalanceNegative bool
	HasOnlyTrial       bool
	HasNewBie          bool
}

type TutorFinance struct {
	Conversion float64
	Count      int64
	BaseCount  int64
	TrialCount int64

	Amount decimal.Decimal
}

type TutorLessons struct {
	LessonsCount int64
	TrialCount   int64
	BaseCount    int64
}
