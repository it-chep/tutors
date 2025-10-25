package dal

import (
	"context"
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

func (r *Repository) GetTutorLessons(ctx context.Context, tutorID int64, from, to time.Time) ([]dto.Lesson, error) {
	sql := `
		select cl.*, s.last_name as last_name, s.first_name as first_name, s.middle_name as middle_name
			from conducted_lessons cl
			join students s on s.id = cl.student_id
		where cl.tutor_id=$1 and cl.created_at between $2 and $3
	`
	var lessons dao.LessonsDAO
	err := pgxscan.Select(ctx, r.pool, &lessons, sql, tutorID, from, to)
	if err != nil {
		return nil, err
	}

	return lessons.ToDomain(ctx), nil
}
