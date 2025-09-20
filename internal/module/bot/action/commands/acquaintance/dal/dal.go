package acquaintance_dal

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

func (d *Dal) ParentOnRegistration(ctx context.Context, parentTG int64) (bool, error) {
	var (
		onRegSql = "select exists (select 1 from registration where tg_id = $1)"
		onReg    bool
	)
	if err := d.pool.QueryRow(ctx, onRegSql, parentTG).Scan(&onReg); err != nil {
		return false, err
	}

	return onReg, nil
}

func (d *Dal) MakeParentKnown(ctx context.Context, parentTG int64, parentFullName string) (bool, error) {
	var (
		setTgSql = "update students set parent_tg_id = $1 where parent_full_name = $2"
	)

	pgTag, err := d.pool.Exec(ctx, setTgSql, parentTG, parentFullName)
	if err != nil {
		return false, err
	}

	defer func() {
		if pgTag.RowsAffected() > 0 {
			_, _ = d.pool.Exec(ctx, "delete from registration where tg_id = $1", parentTG)
		}
	}()

	return pgTag.RowsAffected() > 0, nil
}
