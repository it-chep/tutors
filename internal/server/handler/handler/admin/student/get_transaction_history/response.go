package get_transaction_history

type Transaction struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	IsConfirmed bool   `json:"is_confirmed"`
	Amount      string `json:"amount"`
	IsManual    bool   `json:"is_manual"`
}

type Response struct {
	Transactions         []Transaction `json:"transactions"`
	TransactionsCount    int64         `json:"transactions_count"`
	TotalConfirmedAmount string        `json:"total_confirmed_amount"`
}
