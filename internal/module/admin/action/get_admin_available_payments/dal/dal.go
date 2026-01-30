package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository .
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository .
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// GetAdminPayments получение платежек админа
func (r *Repository) GetAdminPayments(ctx context.Context, adminID int64) ([]dto.Payment, error) {
	sql := `
		select id, bank from payment_cred where admin_id = $1
	`

	var payments dao.Payments
	err := pgxscan.Select(ctx, r.pool, &payments, sql, adminID)
	if err != nil {
		return nil, err
	}

	return payments.ToDomain(), nil
}
