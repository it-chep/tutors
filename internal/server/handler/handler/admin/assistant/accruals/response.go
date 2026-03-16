package accruals

type Accrual struct {
	ID         int64  `json:"id"`
	ActualType string `json:"actual_type"`
	Amount     string `json:"amount"`
	Comment    string `json:"comment"`
	LessonID   int64  `json:"lesson_id"`
	ActualAt   string `json:"actual_at"`
	IsPaid     bool   `json:"is_paid"`
}

type Summary struct {
	Lessons   string `json:"lessons"`
	Penalties string `json:"penalties"`
	Bonuses   string `json:"bonuses"`
	Payable   string `json:"payable"`
	Unpaid    string `json:"unpaid"`
}

type Response struct {
	Accruals []Accrual `json:"accruals"`
	Summary  Summary   `json:"summary"`
}
