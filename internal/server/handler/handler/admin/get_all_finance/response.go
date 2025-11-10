package get_all_finance

type Finance struct {
	Profit   string `json:"profit"`
	CashFlow string `json:"cash_flow"`
}

type Response struct {
	Finance Finance `json:"data"`
}
