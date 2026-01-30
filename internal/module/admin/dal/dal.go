package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// GetStudentsPayments получение платежек студентов
func (r *Repository) GetStudentsPayments(ctx context.Context, studentIDs []int64) (map[int64]dto.Payment, error) {
	sql := `
		select 
		    s.id as "student_id", 
		    pc.id as "payment_id", 
		    pc.bank 
		from students s
		    join payment_cred pc on s.payment_id = pc.id 
		where s.id = any($1)
	`

	var studentPayments []struct {
		StudentID int64  `db:"student_id"`
		PaymentID int64  `db:"payment_id"`
		Bank      string `db:"bank"`
	}

	err := pgxscan.Select(ctx, r.pool, &studentPayments, sql, studentIDs)
	if err != nil {
		return nil, err
	}

	paymentsMap := make(map[int64]dto.Payment, len(studentIDs))
	for _, studentPayment := range studentPayments {
		paymentsMap[studentPayment.StudentID] = dto.Payment{
			ID:   studentPayment.PaymentID,
			Bank: config.Bank(studentPayment.Bank),
		}
	}

	return paymentsMap, nil
}
