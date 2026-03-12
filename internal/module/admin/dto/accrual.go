package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type AccrualActualType int64

const (
	AccrualActualTypeLesson  AccrualActualType = 1
	AccrualActualTypePenalty AccrualActualType = 2
	AccrualActualTypeBonus   AccrualActualType = 3
)

func (t AccrualActualType) String() string {
	switch t {
	case AccrualActualTypeLesson:
		return "lesson"
	case AccrualActualTypePenalty:
		return "penalty"
	case AccrualActualTypeBonus:
		return "bonus"
	default:
		return "unknown"
	}
}

type Accrual struct {
	ID           int64
	TargetUserID int64
	TargetRoleID Role
	ActualTypeID AccrualActualType
	LessonID     int64
	Amount       decimal.Decimal
	Comment      string
	CreatedByID  int64
	ActualAt     time.Time
	IsPaid       bool
}

type AccrualSummary struct {
	Lessons   decimal.Decimal
	Penalties decimal.Decimal
	Bonuses   decimal.Decimal
	Payable   decimal.Decimal
	Unpaid    decimal.Decimal
}
