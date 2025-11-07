package get_tutor_finance

type Finance struct {
	Amount     string  `json:"amount"`
	Wages      string  `json:"wages"` // заработная плата репа
	HoursCount float64 `json:"hours_count"`
}
type Response struct {
	Finance Finance `json:"data"`
}
