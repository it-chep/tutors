package get_all_lessons

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_lessons/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, adminID int64, from, to time.Time) ([]dto.Lesson, error) {
	lessons, err := a.dal.GetAllLessons(ctx, adminID, from, to)
	if err != nil {
		return nil, err
	}

	return lessons, nil
}
