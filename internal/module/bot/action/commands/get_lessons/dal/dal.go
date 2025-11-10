package get_lessons_dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/bot/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"time"
)

type Dal struct {
	pool *pgxpool.Pool
}

func NewDal(pool *pgxpool.Pool) *Dal {
	return &Dal{
		pool: pool,
	}
}

func (d *Dal) GetLessons(ctx context.Context, parentTgID int64) ([]dto.Lesson, error) {
	sql := `
		select cl.* 
		from conducted_lessons cl 
		    join students s
		        on cl.student_id = s.id
		where cl.is_trial is not true and s.parent_tg_id = $1 and cl.created_at between $2 and $3
		order by cl.created_at
	`

	args := []interface{}{
		parentTgID,
		time.Now().Add(-30 * 24 * time.Hour),
		time.Now(),
	}

	var lessons dao.LessonsDAO
	err := pgxscan.Select(ctx, d.pool, &lessons, sql, args...)
	if err != nil {
		return nil, err
	}

	return lessons.ToDomain(), nil
}

func (d *Dal) GetStudentCostByParentTgID(ctx context.Context, parentTgID int64) (decimal.Decimal, error) {
	sql := `select cost_per_hour from students where parent_tg_id = $1`

	var costPerHour decimal.Decimal

	err := pgxscan.Get(ctx, d.pool, &costPerHour, sql, parentTgID)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return costPerHour, nil
}
