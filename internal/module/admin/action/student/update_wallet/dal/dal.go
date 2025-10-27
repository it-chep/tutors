package dal

import (
	"context"

	"github.com/shopspring/decimal"

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

func (r *Repository) UpdateStudentWallet(ctx context.Context, studentID int64, remain decimal.Decimal) error {
	sql := `
		update wallet set balance = $1 where student_id = $2
	`

	args := []interface{}{
		remain,
		studentID,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
