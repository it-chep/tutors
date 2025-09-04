package get_tutor_finance

type Finance struct {
	Conversion int64 `json:"conversion"`
	Count      int64 `json:"count"`
	Amount     int64 `json:"amount"`
}
type Response struct {
	Finance Finance `json:"data"`
}
