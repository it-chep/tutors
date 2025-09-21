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
func (r *Repository) GetFinanceInfo(ctx context.Context, studentID int64, from, to time.Time) (dto.StudentFinance, error) {
	sql := `
		select count(*) as count,
			   sum(
					   case
						   when cl.is_trial = false then (cl.duration_in_minutes / 60.0) * s.cost_per_hour
						   else 0
						   end
			   )        as amount
		from conducted_lessons cl
				 join students s on cl.student_id = s.id
		where cl.student_id = $1
		  and cl.created_at between $2 and $3
	`

	args := []interface{}{
		studentID,
		from,
		to,
	}

	var info dao.StudentFinance
	err := pgxscan.Get(ctx, r.pool, &info, sql, args...)
	if err != nil {
		return dto.StudentFinance{}, err
	}

	return info.ToDomain(), nil
}
