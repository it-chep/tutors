package start_dal

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Dal struct {
	pool *pgxpool.Pool
}

func NewDal(pool *pgxpool.Pool) *Dal {
	return &Dal{
		pool: pool,
	}
}

func (d *Dal) IsKnown(ctx context.Context, parentTG int64) (bool, error) {
	var exists bool
	sql := `
		select exists(select * from students where parent_tg_id = $1)
	`

	if err := d.pool.QueryRow(ctx, sql, parentTG).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
