package update_wallet

import "github.com/shopspring/decimal"

type Request struct {
	Balance string `json:"balance"`
}

func (r Request) BalanceDec() (decimal.Decimal, error) {
	return decimal.NewFromString(r.Balance)
}
