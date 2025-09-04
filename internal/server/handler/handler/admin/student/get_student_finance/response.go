package get_student_finance

type StudentFinance struct {
	Count  int64 `json:"count"`
	Amount int64 `json:"amount"`
}
type Response struct {
	StudentFinance StudentFinance `json:"data"`
}
