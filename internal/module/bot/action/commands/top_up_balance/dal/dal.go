package top_up_balance_dal

import (
	"context"
	"fmt"

	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/samber/lo"

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
			select * from transactions_history 
			where student_id = (select id from students where parent_tg_id = $1) and
				  confirmed_at is null and 
				  amount is null
        `
	)

	err := pgxscan.Get(ctx, d.pool, transaction, sql, parentTG)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Message(ctx, fmt.Sprintf("NIL transaction by parent (parent_tg_id = %d)", parentTG))
			return nil, nil
		}
		logger.Message(ctx, fmt.Sprintf("ERROR TRANSACTIONS (parent_tg_id = %d)", parentTG))
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
		logger.Error(ctx, fmt.Sprintf("InitTransaction ERROR (parent_tg_id = %d, order = %s)", parentTG, order.String()), err)
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

func (d *Dal) SetOrderID(ctx context.Context, transactionID, orderID string, paymentID int64) error {
	sql := `
		update transactions_history 
			set order_id = $2, payment_id = $3
		where id = $1
	`

	_, err := d.pool.Exec(ctx, sql, transactionID, orderID, paymentID)
	return err
}

func (d *Dal) DropTransaction(ctx context.Context, transactionID string) error {
	sql := `
		delete from transactions_history where id = $1
	`

	_, err := d.pool.Exec(ctx, sql, transactionID)
	return err
}

func (d *Dal) AdminIDByParent(ctx context.Context, parentTG int64) (int64, int64, error) {
	var (
		adminID   int64
		paymentID *int64
		sql       = `
			with tutor_id_sel as (
    			select tutor_id, payment_id
					from students 
				where parent_tg_id = $1
			)
			select t.admin_id, s.payment_id
			from tutor_id_sel s
				join tutors t on t.id = s.tutor_id
		`
	)

	if err := d.pool.QueryRow(ctx, sql, parentTG).Scan(&adminID, &paymentID); err != nil {
		return 0, 0, err
	}

	if adminID == 0 {
		return 0, 0, fmt.Errorf("admin id is zero")
	}

	return adminID, lo.FromPtr(paymentID), nil
}
