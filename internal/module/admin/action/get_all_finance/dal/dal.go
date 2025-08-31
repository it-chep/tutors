package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// GetAllFinanceInfo получаем информацию о финансах
func (r *Repository) GetAllFinanceInfo(ctx context.Context, from, to time.Time) (dto.StudentFinance, error) {
	// todo подумать должны ли учитываться неоплаченные занятия ?

	// количество оплаченных занятий - count(*) as count,
	// прибыль - sum(amount) - ЗП каждого репетитора
	// выручка - sum(amount)
	// конверсия - из пробных занятий в оплату
	sql := `
		
	`

	args := []interface{}{
		from,
		to,
	}

	var info dao.StudentFinance
	err := pgxscan.Get(ctx, r.pool, &info, sql, args...)
	if err != nil {
		return dto.StudentFinance{}, err
	}

	return info.ToDomain(), nil
}
