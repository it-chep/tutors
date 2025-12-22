package dto

import (
	"time"
)

type GetAllFinanceDto struct {
	Profit     string
	CashFlow   string
	Debt       string
	TutorsInfo TutorsInfo
}

type Request struct {
	To, From    time.Time
	AdminID     int64
	TgUsernames []string
}

type TutorsInfo struct {
	Hours, Salary string
}
