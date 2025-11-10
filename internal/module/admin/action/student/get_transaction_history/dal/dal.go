package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetTransactionsByRange(ctx context.Context, studentID int64, from, to time.Time) ([]dto.TransactionHistory, error) {
	sql := `
		select * from transactions_history where student_id = $1 and created_at between $2 and $3
	`

	args := []interface{}{
		studentID,
		from,
		to,
	}

	var history dao.TransactionsHistoryDAO
	err := pgxscan.Select(ctx, r.pool, &history, sql, args...)

	return history.ToDomain(), err
}
