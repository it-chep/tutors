package dal

import (
	"context"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

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

// GetExpenses получение расходов на ЗП репетиторов
func (r *Repository) GetExpenses(ctx context.Context, from, to time.Time) (decimal.Decimal, error) {
	sql := `
		SELECT 
			SUM((cl.duration_in_minutes / 60.0) * t.cost_per_hour) as total_tutor_payout
		FROM conducted_lessons cl
		JOIN tutors t ON cl.tutor_id = t.id
		WHERE cl.created_at BETWEEN $1 AND $2
	`

	args := []interface{}{
		from,
		to,
	}

	// todo добавить связь с админом
	var amount pgtype.Numeric
	err := pgxscan.Get(ctx, r.pool, &amount, sql, args...)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return convert.NumericToDecimal(amount), nil
}

// GetTutorsConversion считаем конверсию оплат после триалок у репетиторов
func (r *Repository) GetTutorsConversion(ctx context.Context, from, to time.Time) (float64, error) {
	sql := `
		WITH trial_lessons AS (SELECT cl.student_id,
									  cl.created_at as trial_date,
									  cl.tutor_id
							   FROM conducted_lessons cl
							   WHERE cl.is_trial = true
								 AND cl.created_at BETWEEN $1 AND $2),
			 paid_students AS (SELECT DISTINCT th.student_id
							   FROM transactions_history th
							   WHERE th.confirmed_at BETWEEN $1 AND $2
								 AND th.amount > 0),
			 conversion_data AS (SELECT COUNT(DISTINCT tl.student_id) as total_trial_students,
										COUNT(DISTINCT ps.student_id) as converted_students,
										CASE
											WHEN COUNT(DISTINCT tl.student_id) > 0 THEN
												(COUNT(DISTINCT ps.student_id)::decimal / COUNT(DISTINCT tl.student_id)) * 100
											ELSE 0
											END                       as conversion_rate
								 FROM trial_lessons tl
										  LEFT JOIN paid_students ps ON tl.student_id = ps.student_id)
		SELECT ROUND(conversion_rate, 2) as conversion_rate_percent
		FROM conversion_data
	`

	// todo добавить связь с админом
	args := []interface{}{
		from,
		to,
	}

	var conversion float64
	err := pgxscan.Get(ctx, r.pool, &conversion, sql, args...)
	if err != nil {
		return conversion, err
	}

	return conversion, nil
}

// GetCashFlow получение оборота
func (r *Repository) GetCashFlow(ctx context.Context, from, to time.Time) (decimal.Decimal, error) {
	sql := `
		select sum(amount) from transactions_history where confirmed_at between $1 and $2
	`

	args := []interface{}{
		from,
		to,
	}

	// todo добавить связь с админом
	var amount pgtype.Numeric
	err := pgxscan.Get(ctx, r.pool, &amount, sql, args...)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return convert.NumericToDecimal(amount), nil
}

// GetLessons количество проведенных занятий
func (r *Repository) GetLessons(ctx context.Context, from, to time.Time) (dto.TutorLessons, error) {
	sql := `
		select
			count(*) as lessons_count,
			count(*) filter (where is_trial = false) as base_lessons,
			count(*) filter (where is_trial = true) as trial_lessons
		from conducted_lessons
		where created_at between $1 and $2
	`

	args := []interface{}{
		from,
		to,
	}

	// todo добавить связь с админом
	var lessons dao.TutorLessonsCountDao
	err := pgxscan.Get(ctx, r.pool, &lessons, sql, args...)
	if err != nil {
		return dto.TutorLessons{}, err
	}

	return lessons.ToDomain(), nil
}
