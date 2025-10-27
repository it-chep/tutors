package update_wallet

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/update_wallet/dal"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// Action провести обычное занятие
type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, studentID int64, wallet decimal.Decimal) error {
	return a.dal.UpdateStudentWallet(ctx, studentID, wallet)
}
