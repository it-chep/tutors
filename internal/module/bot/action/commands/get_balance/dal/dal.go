package get_balance_dal

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Dal struct {
	pool *pgxpool.Pool
}

func NewDal(pool *pgxpool.Pool) *Dal {
	return &Dal{
		pool: pool,
	}
}

func (d *Dal) GetBalance(ctx context.Context, parentTG int64) (decimal.Decimal, error) {
	sql := `
		with student as (
    		select id 
				from students 
			where parent_tg_id = $1
		)
		select balance
			from wallet
		where student_id = (select id from student) 
	`
	var balance pgtype.Numeric
	if err := d.pool.QueryRow(ctx, sql, parentTG).Scan(&balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return decimal.Zero, nil
		}
		return decimal.Decimal{}, err
	}

	return convert.NumericToDecimal(balance), nil
}
