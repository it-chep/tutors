package get_admin_available_payments

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_admin_available_payments/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Action .
type Action struct {
	dal *dal.Repository
}

// New .
func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

// Do получение платежек админа
func (a *Action) Do(ctx context.Context, adminID int64) ([]dto.Payment, error) {
	return a.dal.GetAdminPayments(ctx, adminID)
}
