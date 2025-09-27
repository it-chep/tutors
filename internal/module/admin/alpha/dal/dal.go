package alpha_dal

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) UpdateBalance(ctx context.Context, orderNumber string) error {
	sql := `
		with upd as (
			update transactions_history 
				set confirmed_at = now() 
			where id = $1 and confirmed_at is null
			returning student_id, amount
		)
		update wallet w set  
			balance = balance + u.amount
		from upd u
		where w.student_id = u.student_id
	`
	_, err := r.pool.Exec(ctx, sql, orderNumber)
	return err
}
