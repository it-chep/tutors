package alpha_dal

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/bot/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/bot/dto/business"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
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

func (r *Repository) GetOrdersByAmount(ctx context.Context, amount decimal.Decimal) ([]*business.Transaction, error) {
	var (
		daos dao.TransactionDAOs
		sql  = `
			select * from transactions_history
				where amount = $1 and confirmed_at is null
		`
	)

	if err := pgxscan.Select(ctx, r.pool, &daos, sql, amount.Mul(decimal.NewFromInt(100)).String()); err != nil {
		return nil, err
	}

	return daos.ToDomain(), nil
}

func (r *Repository) GetUnconfirmedOrders(ctx context.Context) ([]*business.Transaction, error) {
	var (
		daos dao.TransactionDAOs
		sql  = `
			select * from transactions_history
				where confirmed_at is null and amount is not null and order_id is not null
		`
	)

	if err := pgxscan.Select(ctx, r.pool, &daos, sql); err != nil {
		return nil, err
	}

	return daos.ToDomain(), nil
}
