package dal

import (
	"context"
	"encoding/json"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"strings"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/create_student/dto"
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

// CreateStudent создание студента
func (r *Repository) CreateStudent(ctx context.Context, createDTO dto.CreateRequest) (int64, error) {
	sql := `
		insert into students (first_name, last_name, middle_name, phone, tg, cost_per_hour, subject_id, tutor_id, is_finished_trial, parent_full_name, parent_phone, parent_tg, tg_admin_username, payment_id) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, false, $9, $10, $11, $12, $13)
		returning id
	`
	args := []interface{}{
		strings.TrimSpace(createDTO.FirstName),
		strings.TrimSpace(createDTO.LastName),
		strings.TrimSpace(createDTO.MiddleName),
		createDTO.Phone,
		createDTO.Tg,
		createDTO.CostPerHour,
		createDTO.SubjectID,
		createDTO.TutorID,
		strings.TrimSpace(createDTO.ParentFullName),
		createDTO.ParentPhone,
		createDTO.ParentTg,
		strings.TrimSpace(createDTO.TgAdminUsername),
		createDTO.PaymentID,
	}

	var id int64
	row := r.pool.QueryRow(ctx, sql, args...)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// CreateWallet создание кошелька
func (r *Repository) CreateWallet(ctx context.Context, studentID int64) error {
	sql := `
		insert into wallet ( student_id, balance ) values ( $1, 0 ) 
	`
	args := []interface{}{
		studentID,
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}

// GetDefaultAdminPaymentID получение дефолтной платежки админа
func (r *Repository) GetDefaultAdminPaymentID(ctx context.Context, adminID int64) (int64, error) {
	sql := `
		select id from payment_cred where admin_id = $1 and is_default is true
	`

	var paymentID int64
	err := pgxscan.Get(ctx, r.pool, &paymentID, sql, adminID)
	if err != nil {
		return 0, err
	}

	return paymentID, nil
}

// AddTgToAssistant добавление тг ассистенту
func (r *Repository) AddTgToAssistant(ctx context.Context, assistantID int64, tgAdminUsername string) error {
	sql := `
		update assistant_tgs
		set available_tgs = array(
			select distinct unnest(array_append(available_tgs, $2))
		)
		where user_id = $1
	`
	args := []interface{}{
		assistantID,
		tgAdminUsername,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

// GetPaidFunctions получение оплаченых фич для пользователя
func (r *Repository) GetPaidFunctions(ctx context.Context, adminID int64) (*indto.PaidFunctions, error) {
	sql := "select * from paid_functions where admin_id = $1"
	p := &xo.PaidFunction{}
	if err := pgxscan.Get(ctx, r.pool, p, sql, adminID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	paid := &indto.PaidFunctions{
		AdminID:       p.AdminID,
		PaidFunctions: make(map[string]bool),
	}

	if err := json.Unmarshal(p.Functions, &paid.PaidFunctions); err != nil {
		return nil, err
	}

	return paid, nil
}

// SetUserPaymentUUID сетим пользаку paymentUUID если у админа оплачена эта фича
func (r *Repository) SetUserPaymentUUID(ctx context.Context, studentID int64) error {
	paymentUUID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	sql := `update students set payment_uuid = $1 where id = $2`

	_, err = r.pool.Exec(ctx, sql, paymentUUID, studentID)
	if err != nil {
		return err
	}

	return nil
}
