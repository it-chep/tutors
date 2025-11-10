package get_all_transactions

type Transaction struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	IsConfirmed bool   `json:"is_confirmed"`
	Amount      string `json:"amount"`

	StudentID   int64  `json:"student_id"`
	StudentName string `json:"student_name"`
}

type Response struct {
	Transactions      []Transaction `json:"transactions"`
	TransactionsCount int64         `json:"transactions_count"`
}
