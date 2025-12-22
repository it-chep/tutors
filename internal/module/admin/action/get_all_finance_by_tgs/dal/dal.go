package dal

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance_by_tgs/dto"
	"github.com/samber/lo"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
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

// GetCashFlow получение оборота
func (r *Repository) GetCashFlow(ctx context.Context, req dto.Request, lessons dao.ConductedLessonDAOs) (decimal.Decimal, error) {
	sixty := decimal.NewFromInt(60)
	allMoney := decimal.NewFromFloat(0.0)

	studentsMap := lo.GroupBy(lessons, func(item dao.ConductedLessonDAO) int64 {
		return item.StudentID
	})

	for _, studentInfo := range r.perStudentHours(ctx, req.AdminID) {
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

	return allMoney, nil
}

// GetFinanceInfo получаем прибыль по репетитору
func (r *Repository) GetFinanceInfo(ctx context.Context, req dto.Request, lessons dao.ConductedLessonDAOs) (decimal.Decimal, error) {
	sixty := decimal.NewFromInt(60)
	allMoney := decimal.NewFromFloat(0.0)
	salary := decimal.NewFromFloat(0.0)

	studentsMap := lo.GroupBy(lessons, func(item dao.ConductedLessonDAO) int64 {
		return item.StudentID
	})

	tutorsMap := lo.GroupBy(lessons, func(item dao.ConductedLessonDAO) int64 {
		return item.TutorID
	})

	for _, studentInfo := range r.perStudentHours(ctx, req.AdminID) {
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

	for _, tutorInfo := range r.perTutorHours(ctx, req.AdminID) {
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

func (r *Repository) GetLessonsInfo(ctx context.Context, req dto.Request) (dao.ConductedLessonDAOs, error) {
	sql := `
		select cl.* from conducted_lessons cl
		    	join students s on cl.student_id = s.id
		        join tutors t on cl.tutor_id = t.id
		where t.admin_id = $1
		 	and cl.created_at between $2 and $3
			and cl.is_trial = false
			and s.tg_admin_username = any($4)
	`

	args := []interface{}{
		req.AdminID,
		req.From,
		req.To,
		req.TgUsernames,
	}

	var lessons dao.ConductedLessonDAOs
	err := pgxscan.Select(ctx, r.pool, &lessons, sql, args...)
	if err != nil {
		return lessons, err
	}

	return lessons, nil
}

// GetDebt получение текущей дебиторской задолженности
func (r *Repository) GetDebt(ctx context.Context, req dto.Request) (decimal.Decimal, error) {
	sql := `
		select sum(w.balance) from wallet w 
    		join students s on w.student_id = s.id 
			join tutors t on s.tutor_id = t.id
		where t.admin_id = $1 
		  and balance < 0 
		  and s.tg_admin_username = any($2)
	`

	args := []interface{}{
		req.AdminID,
		req.TgUsernames,
	}

	var debt pgtype.Numeric
	err := pgxscan.Get(ctx, r.pool, &debt, sql, args...)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return convert.NumericToDecimal(debt), nil
}

// GetTutorsInfo получаем информацию по репетиторам
func (r *Repository) GetTutorsInfo(ctx context.Context, req dto.Request, lessons dao.ConductedLessonDAOs) (dto.TutorsInfo, error) {
	sixty := decimal.NewFromInt(60)
	salary := decimal.NewFromFloat(0.0)
	hours := decimal.NewFromFloat(0.0)

	tutorsMap := lo.GroupBy(lessons, func(item dao.ConductedLessonDAO) int64 {
		return item.TutorID
	})

	for _, tutorInfo := range r.perTutorHours(ctx, req.AdminID) {
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

		partOfHour := minutesDecimal.Div(sixty)
		hours = hours.Add(partOfHour)
	}

	return dto.TutorsInfo{
		Hours:  hours.String(),
		Salary: salary.String(),
	}, nil
}
