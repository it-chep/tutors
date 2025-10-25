package get_tutor_lessons

import (
	"context"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/get_tutor_lessons/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, studentID int64, from, to time.Time) (_ []dto.Lesson, err error) {
	return a.dal.GetTutorLessons(ctx, studentID, from, to)
}
