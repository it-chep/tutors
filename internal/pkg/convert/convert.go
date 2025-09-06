package convert

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func NumericToDecimal(numeric pgtype.Numeric) decimal.Decimal {
	if !numeric.Valid {
		return decimal.Zero
	}
	fl, _ := numeric.Float64Value()
	return decimal.NewFromFloat(fl.Float64)
}
