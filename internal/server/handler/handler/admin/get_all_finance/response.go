package get_all_finance

type Finance struct {
	Profit   string `json:"profit"`
	CashFlow string `json:"cash_flow"`

	Conversion        float64 `json:"conversion"`
	LessonsCount      int64   `json:"lessons_count"`
	CountBaseLessons  int64   `json:"base_lessons"`
	CountTrialLessons int64   `json:"trial_lessons"`
}

type Response struct {
	Finance Finance `json:"data"`
}
