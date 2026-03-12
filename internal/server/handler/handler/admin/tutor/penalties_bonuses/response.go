package penalties_bonuses

type Item struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Amount   string `json:"amount"`
	Comment  string `json:"comment"`
	ActualAt string `json:"actual_at"`
	IsPaid   bool   `json:"is_paid"`
}

type Summary struct {
	Penalties string `json:"penalties"`
	Bonuses   string `json:"bonuses"`
	Payable   string `json:"payable"`
	Unpaid    string `json:"unpaid"`
}

type Response struct {
	Items   []Item  `json:"items"`
	Summary Summary `json:"summary"`
}
