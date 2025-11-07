package get_transaction_history

type Transaction struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	IsConfirmed bool   `json:"is_confirmed"`
	Amount      string `json:"amount"`
}

type Response struct {
	Transactions      []Transaction `json:"transactions"`
	TransactionsCount int64         `json:"transactions_count"`
}
