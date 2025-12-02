package top_up_balance_dal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/it-chep/tutors.git/internal/module/bot/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/bot/dto/business"
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

func (d *Dal) TransactionByParent(ctx context.Context, parentTG int64) (*business.Transaction, error) {
	var (
		transaction = &dao.TransactionDAO{}
		sql         = `
            SELECT th.* FROM transactions_history th
            JOIN students s ON th.student_id = s.id
            WHERE s.parent_tg_id = $1 
              AND th.confirmed_at IS NULL 
              AND th.amount IS NULL
            ORDER BY th.created_at DESC
            LIMIT 1
        `
	)

	err := pgxscan.Get(ctx, d.pool, transaction, sql, parentTG)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pending transaction: %w", err)
	}

	return transaction.ToDomain(), nil
}

func (d *Dal) InitTransaction(ctx context.Context, parentTG int64) (string, error) {
	order, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	sql := `
		with student as (
			select id from students where parent_tg_id = $1 
		)
		insert into transactions_history (id, student_id)
			select $2, (select id from student)
	`

	if _, err = d.pool.Exec(ctx, sql, parentTG, order.String()); err != nil {
		return "", err
	}

	return order.String(), nil
}

func (d *Dal) SetTransactionAmount(ctx context.Context, transactionID string, amount int64) error {
	sql := `
		update transactions_history 
			set amount = $2
		where id = $1
	`

	_, err := d.pool.Exec(ctx, sql, transactionID, amount)
	return err
}

func (d *Dal) SetOrderID(ctx context.Context, transactionID, orderID string) error {
	sql := `
		update transactions_history 
			set order_id = $2
		where id = $1
	`

	_, err := d.pool.Exec(ctx, sql, transactionID, orderID)
	return err
}

func (d *Dal) DropTransaction(ctx context.Context, transactionID string) error {
	sql := `
		delete from transactions_history where id = $1
	`

	_, err := d.pool.Exec(ctx, sql, transactionID)
	return err
}

func (d *Dal) AdminIDByParent(ctx context.Context, parentTG int64) (int64, error) {
	var (
		adminID int64
		sql     = `
			with tutor_id_sel as (
    			select tutor_id 
					from students 
				where parent_tg_id = $1
			)
			select admin_id
				from tutors
			where id = (select tutor_id from tutor_id_sel)
		`
	)

	if err := d.pool.QueryRow(ctx, sql, parentTG).Scan(&adminID); err != nil {
		return 0, err
	}

	if adminID == 0 {
		return 0, fmt.Errorf("admin id is zero")
	}

	return adminID, nil
}
