package search_tutor

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/search_tutor/dal"
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

func (a *Action) Do(ctx context.Context, query string) (tutors []dto.Tutor, err error) {
	return a.dal.Search(ctx, query)
}
