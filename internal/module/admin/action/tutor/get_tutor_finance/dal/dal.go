package dal

import (
	"context"
	"github.com/samber/lo"
	"time"

	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/shopspring/decimal"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
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

// GetFinanceInfo получаем прибыль по репетитору
func (r *Repository) GetFinanceInfo(ctx context.Context, tutorID int64, from, to time.Time) (decimal.Decimal, error) {
	lessons, err := r.conductedNotTrialLessons(ctx, tutorID, from, to)
	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	perHours := r.perStudentHours(ctx, tutorID)
	sixty := decimal.NewFromInt(60)
	allMoney := decimal.NewFromFloat(0.0)

	studentsMap := lo.GroupBy(lessons, func(item dao.ConductedLessonDAO) int64 {
		return item.StudentID
	})

	for _, studentInfo := range perHours {
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

	var minutesCount int64
	for _, lesson := range lessons {
		minutesCount += lesson.DurationInMinutes
	}

	minutesDecimal := decimal.NewFromInt(int64(minutesCount))
	tutorCostPerHour := convert.NumericToDecimal(lo.FromPtr(r.perTutorHours(ctx, tutorID)[0].Tutor))

	salary := tutorCostPerHour.Mul(minutesDecimal).Div(sixty)

	amount := allMoney.Add(salary.Mul(decimal.NewFromInt(-1)))

	return amount, nil
}

func (r *Repository) perStudentHours(ctx context.Context, tutorID int64) (moneys []dao.StudentTutorMoney) {
	perHourSQL := `
		select distinct	s.id as student_id,
		    	s.cost_per_hour as student_cost_per_hour
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

func (r *Repository) perTutorHours(ctx context.Context, tutorID int64) (moneys []dao.StudentTutorMoney) {
	perHourSQL := `
		select distinct 
		    t.id as tutor_id, 
		    t.cost_per_hour as tutor_cost_per_hour 
		from tutors t
		where t.id = $1
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

func (r *Repository) GetTutorFinanceInfo(ctx context.Context, tutorID int64, from, to time.Time) (wages decimal.Decimal, hours float64, err error) {
	lessons, err := r.conductedNotTrialLessons(ctx, tutorID, from, to)
	if err != nil {
		return decimal.Zero, 0.0, err
	}

	// ставка репетитора
	tutorHours := r.perTutorHours(ctx, tutorID)[0].Tutor
	tutorCostPerHour := convert.NumericToDecimal(lo.FromPtr(tutorHours))

	var minutesCount int64
	for _, lesson := range lessons {
		minutesCount += lesson.DurationInMinutes
	}

	sixty := decimal.NewFromInt(60)
	minutesDecimal := decimal.NewFromInt(int64(minutesCount))

	wages = tutorCostPerHour.Mul(minutesDecimal).Div(sixty)
	hours = float64(minutesCount) / 60.0

	return wages, hours, nil
}

//// GetLessonsCounters Получение информации об уроках
//func (r *Repository) GetLessonsCounters(ctx context.Context, tutorID int64, from, to time.Time) (dto.TutorLessons, error) {
//	lessons, err := r.conductedNotTrialLessons(ctx, tutorID, from, to)
//	if err != nil {
//		return dto.TutorLessons{}, err
//	}
//
//	counters := dao.TutorLessonsCountDao{
//		LessonsCount: int64(len(lessons)),
//		TrialCount: int64(len(lo.Filter(lessons, func(item dao.ConductedLessonDAO, index int) bool {
//			return item.IsTrial.Bool
//		}))),
//		BaseCount: int64(len(lo.Filter(lessons, func(item dao.ConductedLessonDAO, index int) bool {
//			return !item.IsTrial.Bool
//		}))),
//	}
//	return counters.ToDomain(), nil
//}
