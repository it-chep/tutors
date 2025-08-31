package action

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Aggregator struct {
}

func NewAggregator(pool *pgxpool.Pool) *Aggregator {
	return &Aggregator{}
}
