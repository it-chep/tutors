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

func (r *Repository) GetTransactionsByRange(ctx context.Context, adminID int64, from, to time.Time) (dto.Transactions, error) {
	sql := `
		select th.* 
		from transactions_history th 
			join students s on th.student_id = s.id
			join tutors t on s.tutor_id = t.id
		where t.admin_id = $1 and th.created_at between $2 and $3
	`

	args := []interface{}{
		adminID,
		from,
		to,
	}

	var history dao.TransactionsHistoryDAO
	err := pgxscan.Select(ctx, r.pool, &history, sql, args...)

	return history.ToDomain(), err
}

func (r *Repository) GetStudentsInfo(ctx context.Context, studentIDs []int64) ([]dto.Student, error) {
	sql := `
		select * from students where id = any($1)
	`

	var students dao.StudentsDAO
	err := pgxscan.Select(ctx, r.pool, &students, sql, studentIDs)
	if err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}
