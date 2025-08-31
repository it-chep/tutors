package admin

import (
	"github.com/it-chep/tutors.git/internal/module/admin/action"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Module модуль отвечающий за работу админки
type Module struct {
	Actions *action.Aggregator
}

func New(pool *pgxpool.Pool) *Module {
	actions := action.NewAggregator(pool)

	return &Module{
		Actions: actions,
	}
}
