package get_tutor_finance

type Finance struct {
	Conversion int64  `json:"conversion"`
	Count      int64  `json:"count"`
	Amount     string `json:"amount"`
}
type Response struct {
	Finance Finance `json:"data"`
}
