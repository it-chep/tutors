package dal

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
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

// GetFinanceInfo получаем количество занятий и прибыль по студенту
func (r *Repository) GetFinanceInfo(ctx context.Context, tutorID int64, from, to time.Time) (dto.TutorFinance, error) {
	// todo подумать должны ли учитываться неоплаченные занятия ?
	sql := `
		select count(*) as count, sum(amount) as amount, 123 as conversion -- todo посчитать конверсию
		from transactions_history th 
		    join students s on th.student_id = s.id
		where s.tutor_id = $1 and th.created_at between $2 and $3
	`

	args := []interface{}{
		tutorID,
		from,
		to,
	}

	var info dao.TutorFinance
	err := pgxscan.Get(ctx, r.pool, &info, sql, args...)
	if err != nil {
		return dto.TutorFinance{}, err
	}

	return info.ToDomain(), nil
}
