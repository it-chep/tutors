package get_tutor_finance

type Finance struct {
	Conversion int64  `json:"conversion"`
	Count      int64  `json:"lessons_count"`
	Amount     string `json:"amount"`
	BaseCount  int64  `json:"base_lessons"`
	TrialCount int64  `json:"trial_lessons"`
}
type Response struct {
	Finance Finance `json:"data"`
}
