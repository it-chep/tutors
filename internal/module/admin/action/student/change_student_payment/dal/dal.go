package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
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

// GetStudentAdminID получение студента
func (r *Repository) GetStudentAdminID(ctx context.Context, studentID int64) (int64, error) {
	sql := `
		select t.admin_id from students s join tutors t on s.tutor_id = t.id where s.id = $1
	`

	var studentAdminID int64
	err := pgxscan.Get(ctx, r.pool, &studentAdminID, sql, studentID)

	return studentAdminID, err
}

// ChangePayment изменение платежки по его ID
func (r *Repository) ChangePayment(ctx context.Context, studentID, newPaymentID int64) error {
	sql := `
		update students set payment_id=$1 where id=$2
	`

	_, err := r.pool.Exec(ctx, sql, newPaymentID, studentID)
	return err
}
