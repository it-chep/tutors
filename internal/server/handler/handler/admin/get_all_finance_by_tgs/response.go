package get_all_finance_by_tgs

type Finance struct {
	Profit   string `json:"profit"`
	CashFlow string `json:"cash_flow"`
	Debt     string `json:"debt"`
	Salary   string `json:"salary"`
	Hours    string `json:"hours"`
}

type Response struct {
	Finance Finance `json:"data"`
}
