package alpha_dal

import (
	"context"
	"fmt"

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

func (r *Repository) AdminIDByStudents(ctx context.Context, studentIDs []int64) (map[int64]int64, error) {
	var (
		adminIDByStudent = make(map[int64]int64, len(studentIDs))
		sql              = `
			with tutor_data as (
    			select id as student_id, tutor_id 
					from students 
				where id = any($1)
			)
			select td.student_id, t.admin_id
				from tutor_data td
					join tutors t on t.id = td.tutor_id
		`
	)

	rows, err := r.pool.Query(ctx, sql, studentIDs)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса получения админов по студентам: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var studentID, adminID int64
		if err = rows.Scan(&studentID, &adminID); err != nil {
			return nil, fmt.Errorf("ошибка скана админов и студентов: %s", err.Error())
		}
		adminIDByStudent[studentID] = adminID
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации получения админов и студентов: %s", err.Error())
	}

	return adminIDByStudent, nil
}
