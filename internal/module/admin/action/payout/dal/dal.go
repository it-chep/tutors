package dal

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/it-chep/tutors.git/internal/pkg/transaction/wrapper"
	"github.com/shopspring/decimal"
)

type Repository struct {
	db wrapper.Database
}

type CreatePayoutRequest struct {
	ID           uuid.UUID
	TargetUserID int64
	TargetRoleID int64
	CreatedByID  int64
	Amount       decimal.Decimal
	Comment      string
}

type Payout struct {
	ID                 uuid.UUID       `db:"id"`
	TargetUserID       int64           `db:"target_user_id"`
	TargetRoleID       int64           `db:"target_role_id"`
	CreatedByID        int64           `db:"created_by_id"`
	Amount             decimal.Decimal `db:"amount"`
	Comment            string          `db:"comment"`
	CreatedAt          time.Time       `db:"created_at"`
	ReceiptFileKey     string          `db:"receipt_key"`
	ReceiptFileName    string          `db:"receipt_file_name"`
	ReceiptContentType string          `db:"receipt_content_type"`
	ReceiptUploadedAt  *time.Time      `db:"receipt_uploaded_at"`
	TutorName          string          `db:"tutor_name"`
}

func NewRepository(db wrapper.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, req CreatePayoutRequest) error {
	sql := `
		insert into accrual_payouts (id, target_user_id, target_role_id, created_by_id, amount, comment)
		values ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Pool(ctx).Exec(ctx, sql, req.ID, req.TargetUserID, req.TargetRoleID, req.CreatedByID, req.Amount, req.Comment)
	return err
}

func (r *Repository) MarkTutorAccrualsPaid(ctx context.Context, tutorID int64, payoutID uuid.UUID) (int64, error) {
	sql := `
		update accruals
		set is_paid = true,
			payout_id = $2
		where target_user_id = $1
		  and target_role_id = 3
		  and is_paid = false
	`

	tag, err := r.db.Pool(ctx).Exec(ctx, sql, tutorID, payoutID)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (r *Repository) GetNearestPayoutWithoutReceipt(ctx context.Context, tutorID int64) (Payout, error) {
	sql := `
		select
			ap.id,
			ap.target_user_id,
			ap.target_role_id,
			ap.created_by_id,
			ap.amount,
			coalesce(ap.comment, '') as comment,
			ap.created_at,
			coalesce(ap.receipt_key, '') as receipt_key,
			coalesce(ap.receipt_file_name, '') as receipt_file_name,
			coalesce(ap.receipt_content_type, '') as receipt_content_type,
			ap.receipt_uploaded_at,
			u.full_name as tutor_name
		from accrual_payouts ap
		join users u on u.tutor_id = ap.target_user_id
		where ap.target_user_id = $1
		  and ap.target_role_id = 3
		  and ap.receipt_key is null
		order by ap.created_at asc
		limit 1
	`

	var payout Payout
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &payout, sql, tutorID)
	return payout, err
}

func (r *Repository) AttachReceipt(ctx context.Context, payoutID uuid.UUID, key, fileName, contentType string) error {
	sql := `
		update accrual_payouts
		set receipt_key = $2,
			receipt_file_name = $3,
			receipt_content_type = $4,
			receipt_uploaded_at = now()
		where id = $1
	`

	_, err := r.db.Pool(ctx).Exec(ctx, sql, payoutID, key, fileName, contentType)
	return err
}

func (r *Repository) GetPayout(ctx context.Context, payoutID uuid.UUID) (Payout, error) {
	sql := `
		select
			ap.id,
			ap.target_user_id,
			ap.target_role_id,
			ap.created_by_id,
			ap.amount,
			coalesce(ap.comment, '') as comment,
			ap.created_at,
			coalesce(ap.receipt_key, '') as receipt_key,
			coalesce(ap.receipt_file_name, '') as receipt_file_name,
			coalesce(ap.receipt_content_type, '') as receipt_content_type,
			ap.receipt_uploaded_at,
			u.full_name as tutor_name
		from accrual_payouts ap
		join users u on u.tutor_id = ap.target_user_id
		where ap.id = $1
	`

	var payout Payout
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &payout, sql, payoutID)
	return payout, err
}

func (r *Repository) ListTutorPayouts(ctx context.Context, tutorID int64, from, to time.Time) ([]Payout, error) {
	return r.selectPayouts(ctx, `
		where ap.target_user_id = $1
		  and ap.target_role_id = 3
		  and ap.created_at between $2 and $3
		order by ap.created_at desc, ap.id desc
	`, tutorID, from, to)
}

func (r *Repository) ListAllPayouts(ctx context.Context, from, to time.Time) ([]Payout, error) {
	return r.selectPayouts(ctx, `
		where ap.target_role_id = 3
		  and ap.created_at between $1 and $2
		order by ap.created_at desc, ap.id desc
	`, from, to)
}

func (r *Repository) ListAdminPayouts(ctx context.Context, adminID int64, from, to time.Time) ([]Payout, error) {
	return r.selectPayouts(ctx, `
		join tutors t on t.id = ap.target_user_id
		where ap.target_role_id = 3
		  and t.admin_id = $1
		  and ap.created_at between $2 and $3
		order by ap.created_at desc, ap.id desc
	`, adminID, from, to)
}

func (r *Repository) ListAssistantPayouts(ctx context.Context, assistantID, adminID int64, from, to time.Time) ([]Payout, error) {
	return r.selectPayouts(ctx, `
		join tutors t on t.id = ap.target_user_id
		left join assistant_tgs at on at.user_id = $1
		where ap.target_role_id = 3
		  and t.admin_id = $2
		  and ap.created_at between $3 and $4
		  and (
			at.available_tg_ids is null
			or array_length(at.available_tg_ids, 1) = 0
			or t.tg_admin_username_id = any(at.available_tg_ids)
		  )
		order by ap.created_at desc, ap.id desc
	`, assistantID, adminID, from, to)
}

func (r *Repository) CanAdminManageTutor(ctx context.Context, tutorID, adminID int64) (bool, error) {
	sql := `select exists(select 1 from tutors where id = $1 and admin_id = $2)`

	var exists bool
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &exists, sql, tutorID, adminID)
	return exists, err
}

func (r *Repository) CanAssistantAccessTutor(ctx context.Context, tutorID, assistantID, adminID int64) (bool, error) {
	sql := `
		select exists(
			select 1
			from tutors t
			left join assistant_tgs at on at.user_id = $2
			where t.id = $1
			  and t.admin_id = $3
			  and (
				at.available_tg_ids is null
				or array_length(at.available_tg_ids, 1) = 0
				or t.tg_admin_username_id = any(at.available_tg_ids)
			  )
		)
	`

	var exists bool
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &exists, sql, tutorID, assistantID, adminID)
	return exists, err
}

func (r *Repository) selectPayouts(ctx context.Context, condition string, args ...any) ([]Payout, error) {
	sql := `
		select
			ap.id,
			ap.target_user_id,
			ap.target_role_id,
			ap.created_by_id,
			ap.amount,
			coalesce(ap.comment, '') as comment,
			ap.created_at,
			coalesce(ap.receipt_key, '') as receipt_key,
			coalesce(ap.receipt_file_name, '') as receipt_file_name,
			coalesce(ap.receipt_content_type, '') as receipt_content_type,
			ap.receipt_uploaded_at,
			u.full_name as tutor_name
		from accrual_payouts ap
		join users u on u.tutor_id = ap.target_user_id
	` + condition

	var payouts []Payout
	if err := pgxscan.Select(ctx, r.db.Pool(ctx), &payouts, sql, args...); err != nil {
		return nil, err
	}

	return payouts, nil
}
