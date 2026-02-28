package dal

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/transaction/wrapper"
	"github.com/shopspring/decimal"
)

type Repository struct {
	db wrapper.Database
}

func NewRepository(db wrapper.Database) *Repository {
	return &Repository{
		db: db,
	}
}

// GetStudentAdminID получение admin_id по студенту
func (r *Repository) GetStudentAdminID(ctx context.Context, studentID int64) (int64, error) {
	sql := `select t.admin_id from students s join tutors t on s.tutor_id = t.id where s.id = $1`
	var adminID int64
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &adminID, sql, studentID)
	return adminID, err
}

// AddManualTransaction создание ручной транзакции (сразу подтверждённой)
func (r *Repository) AddManualTransaction(ctx context.Context, studentID int64, amount int64) (uuid.UUID, error) {
	sql := `
		insert into transactions_history (id, student_id, amount, created_at, confirmed_at, is_manual)
		values ($1, $2, $3, $4, $4, true)
		returning id
	`
	id := uuid.New()
	now := time.Now()
	_, err := r.db.Pool(ctx).Exec(ctx, sql, id, studentID, amount, now)
	return id, err
}

// GetStudentWallet получение кошелька студента
func (r *Repository) GetStudentWallet(ctx context.Context, studentID int64) (dto.Wallet, error) {
	sql := `
		select * from wallet where student_id = $1
	`

	var wallet dao.Wallet
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &wallet, sql, studentID)
	if err != nil {
		return dto.Wallet{}, err
	}
	return wallet.ToDomain(), err
}

// UpdateStudentWallet создание ручной транзакции (сразу подтверждённой)
func (r *Repository) UpdateStudentWallet(ctx context.Context, studentID int64, remain decimal.Decimal) error {
	sql := `
		update wallet set balance = $1 where student_id = $2
	`

	args := []interface{}{
		remain,
		studentID,
	}

	_, err := r.db.Pool(ctx).Exec(ctx, sql, args...)
	return err
}
