package dal

import (
	"context"
	"github.com/samber/lo"
	"time"

	"github.com/it-chep/tutors.git/internal/pkg/convert"
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

// GetLessonsCounters Получение информации об уроках
func (r *Repository) GetLessonsCounters(ctx context.Context, tutorID int64, from, to time.Time) (dto.TutorLessons, error) {
	lessons, err := r.conductedNotTrialLessons(ctx, tutorID, from, to)
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
func (r *Repository) GetFinanceInfo(ctx context.Context, tutorID int64, from, to time.Time) (decimal.Decimal, error) {
	lessons, err := r.conductedNotTrialLessons(ctx, tutorID, from, to)
	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	perHours := r.perHours(ctx, tutorID)
	sixty := decimal.NewFromInt(60)
	allMoney := decimal.NewFromFloat(0.0)

	studentsMap := lo.GroupBy(lessons, func(item dao.ConductedLessonDAO) int64 {
		return item.StudentID
	})

	for _, studentInfo := range perHours {
		studentLessons, ok := studentsMap[studentInfo.StudentID]
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

	var minutesCount int64
	for _, lesson := range lessons {
		minutesCount += lesson.DurationInMinutes
	}

	minutesDecimal := decimal.NewFromInt(int64(minutesCount))
	tutorCostPerHour := convert.NumericToDecimal(lo.FromPtr(perHours[0].Tutor))

	salary := tutorCostPerHour.Mul(minutesDecimal).Div(sixty)

	amount := allMoney.Add(salary.Mul(decimal.NewFromInt(-1)))

	return amount, nil
}

func (r *Repository) perHours(ctx context.Context, tutorID int64) (moneys []dao.StudentTutorMoney) {
	perHourSQL := `
		select 	s.id as student_id,
		        s.tutor_id as tutor_id,
		    	s.cost_per_hour as student_cost_per_hour,
		       	t.cost_per_hour as tutor_cost_per_hour
		from students s
			join tutors t on s.tutor_id = t.id
		where s.tutor_id = $1
	`

	err := pgxscan.Select(ctx, r.pool, &moneys, perHourSQL, tutorID)
	if err != nil {
		return nil
	}

	return
}

func (r *Repository) conductedNotTrialLessons(ctx context.Context, tutorID int64, from, to time.Time) (dao.ConductedLessonDAOs, error) {
	sql := `
		select * from conducted_lessons cl
		where cl.tutor_id = $1
		 	and cl.created_at between $2 and $3
			and cl.is_trial = false
	`

	args := []interface{}{
		tutorID,
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
