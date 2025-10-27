package dal

import (
	"context"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
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
	lessons, err := r.conductedNotTrialLessons(ctx, studentID, from, to)
	if err != nil {
		return dto.StudentFinance{}, err
	}

	perHours := r.perHours(ctx, studentID)

	var minutesCount int64
	for _, lesson := range lessons {
		minutesCount += lesson.DurationInMinutes
	}

	studentCostPerHour := convert.NumericToDecimal(lo.FromPtr(perHours.Student))
	tutorCostPerHour := convert.NumericToDecimal(lo.FromPtr(perHours.Tutor))

	minutesDecimal := decimal.NewFromInt(int64(minutesCount))
	sixty := decimal.NewFromInt(60)

	allMoney := studentCostPerHour.Mul(minutesDecimal).Div(sixty)
	salary := tutorCostPerHour.Mul(minutesDecimal).Div(sixty)

	amount := allMoney.Add(salary.Mul(decimal.NewFromInt(-1)))
	return dto.StudentFinance{
		Count:  int64(len(lessons)),
		Amount: amount,
	}, nil
}

func (r *Repository) perHours(ctx context.Context, studentID int64) (money dao.StudentTutorMoney) {
	perHourSQL := `
		select distinct s.id as student_id,
		        s.tutor_id as tutor_id,
		    	s.cost_per_hour as student_cost_per_hour,
		       	t.cost_per_hour as tutor_cost_per_hour
		from students s
			join tutors t on s.tutor_id = t.id
		where s.id = $1
	`

	err := pgxscan.Get(ctx, r.pool, &money, perHourSQL, studentID)
	if err != nil {
		return dao.StudentTutorMoney{}
	}

	return
}

func (r *Repository) conductedNotTrialLessons(ctx context.Context, studentID int64, from, to time.Time) (dao.ConductedLessonDAOs, error) {
	sql := `
		select * from conducted_lessons cl
		where cl.student_id = $1
		  	and cl.created_at between $2 and $3
			and cl.is_trial = false
	`

	args := []interface{}{
		studentID,
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
