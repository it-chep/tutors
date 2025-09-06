package get_all_finance

type Finance struct {
	Profit   string `json:"profit"`
	CashFlow string `json:"cash_flow"`

	Conversion   float64 `json:"conversion"`
	LessonsCount int64   `json:"lessons_count"`
}

type Response struct {
	Finance Finance `json:"data"`
}
