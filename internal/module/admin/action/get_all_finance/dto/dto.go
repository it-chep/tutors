package dto

type GetAllFinanceDto struct {
	Profit     string
	CashFlow   string
	Debt       string
	TutorsInfo TutorsInfo
}

type TutorsInfo struct {
	Hours, Salary string
}
