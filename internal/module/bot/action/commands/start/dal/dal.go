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
		with known as (
    		select exists (select 1 from students where parent_tg_id = $1) as is_known
		), ins as (
     		insert into registration (tg_id) select $1
        		where (select not is_known from known)
     		on conflict (tg_id) do nothing
		)
		select is_known from known
	`

	if err := d.pool.QueryRow(ctx, sql, parentTG).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
