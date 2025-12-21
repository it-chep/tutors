package dal

import (
	"context"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
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

func (r *Repository) GetStudentLessons(ctx context.Context, studentID int64, from, to time.Time) ([]dto.Lesson, error) {
	sql := `
		select cl.*, s.last_name as last_name, s.first_name as first_name, s.middle_name as middle_name
			from conducted_lessons cl
			join students s on s.id = cl.student_id
		where cl.student_id=$1 and cl.created_at between $2 and $3
	`

	args := []interface{}{studentID, from, to}

	if dto.IsAssistantRole(ctx) {
		sql = `
			select cl.*, s.last_name as last_name, s.first_name as first_name, s.middle_name as middle_name
			from conducted_lessons cl
			        join students s on s.id = cl.student_id
			where cl.student_id=$1 and cl.created_at between $2 and $3
				and (
			-- Условие A: Если у ассистента есть TG, используем их
		      s.tg_admin_username = any(
				  SELECT available_tgs
				  FROM assistant_tgs 
				  WHERE user_id = $4
					AND available_tgs IS NOT NULL 
					AND array_length(available_tgs, 1) > 0
			  )
			  -- Условие B: Если у ассистента нет TG (пустой массив или нет записи), показываем всех
			  OR NOT EXISTS (
				  SELECT 1
				  FROM assistant_tgs 
				  WHERE user_id = $4
					AND available_tgs IS NOT NULL 
					AND array_length(available_tgs, 1) > 0
			  )
			)
		`
		args = append(args, userCtx.UserIDFromContext(ctx))
	}

	var lessons dao.LessonsDAO
	err := pgxscan.Select(ctx, r.pool, &lessons, sql, args...)
	if err != nil {
		return nil, err
	}

	return lessons.ToDomain(ctx), nil
}
