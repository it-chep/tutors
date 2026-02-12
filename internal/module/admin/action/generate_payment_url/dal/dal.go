package dal

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// GetStudentByIDAndPaymentUUID получаем студента по ID и UUID
func (r *Repository) GetStudentByIDAndPaymentUUID(ctx context.Context, studentID int64, studentUUID string) (dto.Student, error) {
	sql := `
		select * from students where id = $1 and payment_uuid = $2
	`

	var student dao.StudentDAO
	err := pgxscan.Get(ctx, r.pool, &student, sql, studentID, studentUUID)
	if err != nil {
		return dto.Student{}, err
	}

	return student.ToDomain(), nil
}

// CountLastMinuteTransactions сколько транзакций у студента за последнюю минуту
func (r *Repository) CountLastMinuteTransactions(ctx context.Context, studentID int64) (int64, error) {
	sql := `
		select count(*) 
		from transactions_history 
		where student_id = $1 
		  and created_at between now() - interval '1 minute' and now()
	  `

	var count int64
	err := pgxscan.Get(ctx, r.pool, &count, sql, studentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

// InitTransaction создаем транзакцию
func (r *Repository) InitTransaction(ctx context.Context, studentID, amount, paymentID int64) (string, error) {
	internalTransactionUUID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	sql := `
		insert into transactions_history (id, student_id, amount, payment_id) 
		values ($1, $2, $3, $4)
	`

	args := []interface{}{
		internalTransactionUUID.String(),
		studentID,
		amount,
		paymentID,
	}

	if _, err = r.pool.Exec(ctx, sql, args...); err != nil {
		return "", err
	}

	return internalTransactionUUID.String(), nil
}

// SetTransactionOrder обновляем транзакции
func (r *Repository) SetTransactionOrder(ctx context.Context, internalTransactionUUID, orderID string) error {
	sql := `
		update transactions_history set order_id = $1 where id = $2
	`

	args := []interface{}{
		orderID,
		internalTransactionUUID,
	}

	if _, err := r.pool.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

// PaymentAndAdminByStudent получение adminID и paymentID по студенту
func (r *Repository) PaymentAndAdminByStudent(ctx context.Context, studentID int64) (int64, int64, error) {
	var (
		adminID   int64
		paymentID *int64
		sql       = `
			select t.admin_id, s.payment_id 
			from students s 
				join tutors t on s.tutor_id = t.id
			where s.id = $1
		`
	)

	if err := r.pool.QueryRow(ctx, sql, studentID).Scan(&adminID, &paymentID); err != nil {
		return 0, 0, err
	}

	if adminID == 0 {
		return 0, 0, fmt.Errorf("admin id is zero")
	}

	return adminID, lo.FromPtr(paymentID), nil
}

// DropTransaction удаление транзакций
func (r *Repository) DropTransaction(ctx context.Context, transactionID string) error {
	sql := `
		delete from transactions_history where id = $1
	`

	_, err := r.pool.Exec(ctx, sql, transactionID)
	return err
}

func (d *Repository) PhoneByStudent(ctx context.Context, studentID int64) (string, error) {
	var (
		phone string
		sql   = `
			select parent_phone from students where id = $1 
        `
	)

	if err := d.pool.QueryRow(ctx, sql, studentID).Scan(&phone); err != nil {
		return "", fmt.Errorf("failed to get phone by student: %w", err)
	}

	return phone, nil
}
