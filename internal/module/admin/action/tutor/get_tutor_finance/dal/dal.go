package dal

import (
	"context"
	"time"

	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

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

func (r *Repository) GetConversion(ctx context.Context, tutorID int64, from, to time.Time) (float64, error) {
	sql := `
		WITH trial_lessons AS (SELECT cl.student_id,
							  cl.created_at as trial_date,
							  cl.tutor_id
					   FROM conducted_lessons cl
					   WHERE cl.is_trial = true
						 AND cl.created_at BETWEEN $1 AND $2
						 and cl.tutor_id = $3),
		 paid_students AS (SELECT DISTINCT th.student_id
						   FROM transactions_history th
									join students s on th.student_id = s.id
						   WHERE th.confirmed_at BETWEEN $1 AND $2
							 and s.tutor_id = $3
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

	args := []interface{}{
		from,
		to,
		tutorID,
	}

	var conversion float64
	err := pgxscan.Get(ctx, r.pool, &conversion, sql, args...)
	return conversion, err
}

// GetLessons Получение информации об уроках
func (r *Repository) GetLessons(ctx context.Context, tutorID int64, from, to time.Time) (dto.TutorLessons, error) {
	sql := `
		select
			count(*) as lessons_count,
			count(*) filter (where is_trial = false) as base_lessons,
			count(*) filter (where is_trial = true) as trial_lessons
		from conducted_lessons
		where tutor_id = $1 and created_at between $2 and $3
	`

	var lessons dao.TutorLessonsCountDao
	err := pgxscan.Get(ctx, r.pool, &lessons, sql, tutorID, from, to)
	if err != nil {
		return dto.TutorLessons{}, err
	}

	return lessons.ToDomain(), nil
}

// GetFinanceInfo получаем прибыль по студенту
func (r *Repository) GetFinanceInfo(ctx context.Context, tutorID int64, from, to time.Time) (decimal.Decimal, error) {
	sql := `
		select sum(amount) as amount
		from transactions_history th 
		    join students s on th.student_id = s.id
		where s.tutor_id = $1 and th.confirmed_at between $2 and $3
	`

	args := []interface{}{
		tutorID,
		from,
		to,
	}

	var amount pgtype.Numeric
	err := pgxscan.Get(ctx, r.pool, &amount, sql, args...)
	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	return convert.NumericToDecimal(amount), nil
}
