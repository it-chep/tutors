package get_student_finance

type StudentFinance struct {
	Count                int64  `json:"count"`
	Amount               string `json:"amount"`
	TotalConfirmedAmount string `json:"total_confirmed_amount"`
}
type Response struct {
	StudentFinance StudentFinance `json:"data"`
}
