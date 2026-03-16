package payouts

type Request struct {
	Amount  int64  `json:"amount"`
	Comment string `json:"comment"`
}
