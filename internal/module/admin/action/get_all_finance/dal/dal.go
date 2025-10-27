package dal

import (
	"context"
	"github.com/samber/lo"
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
func (r *Repository) GetExpenses(ctx context.Context, from, to time.Time, adminID int64) (decimal.Decimal, error) {
	sql := `
		select 
			sum((cl.duration_in_minutes / 60.0) * t.cost_per_hour) as total_tutor_payout
		from conducted_lessons cl
		join tutors t on cl.tutor_id = t.id
		where t.admin_id = $3 and cl.created_at between $1 and $2
	`

	args := []interface{}{
		from,
		to,
		adminID,
	}

	var amount pgtype.Numeric
	err := pgxscan.Get(ctx, r.pool, &amount, sql, args...)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return convert.NumericToDecimal(amount), nil
}

// GetTutorsConversion считаем конверсию оплат после триалок у репетиторов
func (r *Repository) GetTutorsConversion(ctx context.Context, from, to time.Time, adminID int64) (float64, error) {
	sql := `
		WITH trial_lessons AS (SELECT cl.student_id,
									  cl.created_at as trial_date,
									  cl.tutor_id
							   FROM conducted_lessons cl
										join tutors t on cl.tutor_id = t.id
							   WHERE t.admin_id = $3 and cl.is_trial = true
								 AND cl.created_at BETWEEN $1 AND $2),
			 paid_students AS (SELECT DISTINCT th.student_id
							   FROM transactions_history th
									join students s on th.student_id = s.id
									join tutors t on s.tutor_id = t.id
							   WHERE t.admin_id = $3 and th.confirmed_at BETWEEN $1 AND $2
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
		adminID,
	}

	var conversion float64
	err := pgxscan.Get(ctx, r.pool, &conversion, sql, args...)
	if err != nil {
		return conversion, err
	}

	return conversion, nil
}

// GetCashFlow получение оборота
func (r *Repository) GetCashFlow(ctx context.Context, from, to time.Time, adminID int64) (decimal.Decimal, error) {
	sql := `
		select sum(th.amount) 
		from transactions_history th 
			join students s on th.student_id = s.id
			join tutors t on s.tutor_id = t.id
		where t.admin_id = $3 and th.confirmed_at between $1 and $2
	`

	args := []interface{}{
		from,
		to,
		adminID,
	}

	var amount pgtype.Numeric
	err := pgxscan.Get(ctx, r.pool, &amount, sql, args...)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return convert.NumericToDecimal(amount), nil
}

// GetLessons количество проведенных занятий
func (r *Repository) GetLessons(ctx context.Context, from, to time.Time, adminID int64) (dto.TutorLessons, error) {
	sql := `
		select cl.* from conducted_lessons cl
		join tutors t on cl.tutor_id = t.id
		where t.admin_id = $3 and created_at between $1 and $2
	`

	args := []interface{}{
		from,
		to,
		adminID,
	}

	var lessons []dao.ConductedLessonDAO
	err := pgxscan.Get(ctx, r.pool, &lessons, sql, args...)
	if err != nil {
		return dto.TutorLessons{}, err
	}

	counters := dao.TutorLessonsCountDao{
		LessonsCount: int64(len(lessons)),
		TrialCount: int64(len(lo.Filter(lessons, func(item dao.ConductedLessonDAO, index int) bool {
			return item.IsTrial.Bool
		}))),
		BaseCount: int64(len(lo.Filter(lessons, func(item dao.ConductedLessonDAO, index int) bool {
			return !item.IsTrial.Bool
		}))),
	}
	return counters.ToDomain(), nil
}

// GetFinanceInfo получаем прибыль по репетитору
func (r *Repository) GetFinanceInfo(ctx context.Context, from, to time.Time, adminID int64) (decimal.Decimal, error) {
	lessons, err := r.conductedNotTrialLessons(ctx, adminID, from, to)
	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	sixty := decimal.NewFromInt(60)
	allMoney := decimal.NewFromFloat(0.0)
	salary := decimal.NewFromFloat(0.0)

	studentsMap := lo.GroupBy(lessons, func(item dao.ConductedLessonDAO) int64 {
		return item.StudentID
	})

	tutorsMap := lo.GroupBy(lessons, func(item dao.ConductedLessonDAO) int64 {
		return item.TutorID
	})

	for _, studentInfo := range r.perStudentHours(ctx, adminID) {
		studentLessons, ok := studentsMap[lo.FromPtr(studentInfo.StudentID)]
		if !ok {
			continue
		}
		// сколько отзанимался типочек
		var minutesCount int64
		for _, studentLesson := range studentLessons {
			minutesCount += studentLesson.DurationInMinutes
		}

		// считаем сколько он занес
		studentCostPerHour := convert.NumericToDecimal(lo.FromPtr(studentInfo.Student))
		minutesDecimal := decimal.NewFromInt(int64(minutesCount))
		userMoney := studentCostPerHour.Mul(minutesDecimal).Div(sixty)

		allMoney = allMoney.Add(userMoney)
	}

	for _, tutorInfo := range r.perTutorHours(ctx, adminID) {
		tutorsLessons, ok := tutorsMap[lo.FromPtr(tutorInfo.TutorID)]
		if !ok {
			continue
		}

		// сколько отвел занятий типочек
		var minutesCount int64
		for _, tutorLesson := range tutorsLessons {
			minutesCount += tutorLesson.DurationInMinutes
		}

		// считаем сколько на него надо потратить
		tutorCostPerHour := convert.NumericToDecimal(lo.FromPtr(tutorInfo.Tutor))
		minutesDecimal := decimal.NewFromInt(int64(minutesCount))
		tutorMoney := tutorCostPerHour.Mul(minutesDecimal).Div(sixty)

		salary = salary.Add(tutorMoney)
	}

	amount := allMoney.Add(salary.Mul(decimal.NewFromInt(-1)))

	return amount, nil
}

func (r *Repository) perStudentHours(ctx context.Context, adminID int64) (moneys []dao.StudentTutorMoney) {
	perHourSQL := `
		select distinct s.id            as student_id,
			   s.cost_per_hour as student_cost_per_hour
		from students s
				 join tutors t on s.tutor_id = t.id
		where t.admin_id = $1
	`

	err := pgxscan.Select(ctx, r.pool, &moneys, perHourSQL, adminID)
	if err != nil {
		return nil
	}

	return
}

func (r *Repository) perTutorHours(ctx context.Context, adminID int64) (moneys []dao.StudentTutorMoney) {
	perHourSQL := `
		select distinct s.tutor_id      as tutor_id,
			   t.cost_per_hour as tutor_cost_per_hour
		from students s
				 join tutors t on s.tutor_id = t.id
		where t.admin_id = $1
	`

	err := pgxscan.Select(ctx, r.pool, &moneys, perHourSQL, adminID)
	if err != nil {
		return nil
	}

	return
}

func (r *Repository) conductedNotTrialLessons(ctx context.Context, adminID int64, from, to time.Time) (dao.ConductedLessonDAOs, error) {
	sql := `
		select cl.* from conducted_lessons cl
		         join tutors t on cl.tutor_id = t.id
		where t.admin_id = $1
		 	and cl.created_at between $2 and $3
			and cl.is_trial = false
	`

	args := []interface{}{
		adminID,
		from,
		to,
	}

	var lessons dao.ConductedLessonDAOs
	err := pgxscan.Select(ctx, r.pool, &lessons, sql, args...)
	if err != nil {
		return lessons, err
	}

	return lessons, nil
}
