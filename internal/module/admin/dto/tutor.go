package dto

import "github.com/shopspring/decimal"

type Tutor struct {
	ID          int64
	FullName    string
	Phone       string
	Tg          string
	CostPerHour string
	SubjectID   int64
	AdminID     int64
}

type TutorFinance struct {
	Conversion int64
	Count      int64
	Amount     decimal.Decimal
}
