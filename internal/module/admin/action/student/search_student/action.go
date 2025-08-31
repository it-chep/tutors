package search_student

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/search_student/dal"
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

func (a *Action) Do(ctx context.Context, query string) (students []dto.Student, err error) {
	return a.dal.SearchStudent(ctx, query)
}
