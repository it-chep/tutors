package dal

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

// ChangeAllPayments меняет платежку у всех студентов данного администратора
func (r *Repository) ChangeAllPayments(ctx context.Context, adminID, newPaymentID int64) error {
	sql := `
		update students set payment_id = $1
			where tutor_id in (select id from tutors where admin_id = $2)
	`
	_, err := r.pool.Exec(ctx, sql, newPaymentID, adminID)
	return err
}
