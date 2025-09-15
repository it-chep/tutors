package dto

type GetAllFinanceDto struct {
	Profit     string
	CashFlow   string
	Conversion float64

	CountLessons      int64
	CountBaseLessons  int64
	CountTrialLessons int64
}
