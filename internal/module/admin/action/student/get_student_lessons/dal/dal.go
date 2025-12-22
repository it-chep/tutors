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
                    not exists (
                        select 1 
                        from assistant_tgs at
                        where at.user_id = $4
                          and at.available_tgs is not null
                          and array_length(at.available_tgs, 1) > 0
                    )
                    or s.tg_admin_username in (
                        select unnest(at.available_tgs)
                        from assistant_tgs at
                        where at.user_id = $4
                          and at.available_tgs is not null
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
