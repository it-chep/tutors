package dal

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/transaction/wrapper"
	"github.com/shopspring/decimal"
)

type Repository struct {
	db wrapper.Database
}

type CreateAccrualRequest struct {
	TargetUserID int64
	TargetRoleID int64
	ActualTypeID int64
	Amount       decimal.Decimal
	Comment      string
	CreatedByID  int64
	ActualAt     time.Time
}

type accrualDAO struct {
	ID           int64           `db:"id"`
	TargetUserID int64           `db:"target_user_id"`
	TargetRoleID int64           `db:"target_role_id"`
	ActualTypeID int64           `db:"actual_type_id"`
	LessonID     *int64          `db:"lesson_id"`
	Amount       decimal.Decimal `db:"amount"`
	Comment      string          `db:"comment"`
	CreatedByID  *int64          `db:"created_by_id"`
	ActualAt     time.Time       `db:"actual_at"`
	IsPaid       bool            `db:"is_paid"`
}

func NewRepository(db wrapper.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, req CreateAccrualRequest) error {
	sql := `
		insert into accruals (
			target_user_id, target_role_id, actual_type_id, amount, comment, created_by_id, actual_at
		)
		values ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Pool(ctx).Exec(ctx, sql, req.TargetUserID, req.TargetRoleID, req.ActualTypeID, req.Amount, req.Comment, req.CreatedByID, req.ActualAt)
	return err
}

func (r *Repository) UpsertLessonAccrual(ctx context.Context, lessonID, tutorID int64, amount decimal.Decimal, actualAt time.Time) error {
	sql := `
		insert into accruals (
			target_user_id, target_role_id, actual_type_id, lesson_id, amount, actual_at
		)
		values ($1, $2, $3, $4, $5, $6)
		on conflict (lesson_id) do update
		set amount = excluded.amount,
			actual_at = excluded.actual_at
	`
	_, err := r.db.Pool(ctx).Exec(ctx, sql, tutorID, int64(dto.TutorRole), int64(dto.AccrualActualTypeLesson), lessonID, amount, actualAt.UTC())
	return err
}

func (r *Repository) DeleteLessonAccrual(ctx context.Context, lessonID int64) error {
	sql := `delete from accruals where lesson_id = $1`
	_, err := r.db.Pool(ctx).Exec(ctx, sql, lessonID)
	return err
}

func (r *Repository) List(ctx context.Context, targetUserID, targetRoleID int64, from, to time.Time) ([]dto.Accrual, error) {
	sql := `
		select id, target_user_id, target_role_id, actual_type_id, lesson_id, amount, coalesce(comment, '') as comment, created_by_id, actual_at, is_paid
		from accruals
		where target_user_id = $1
		  and target_role_id = $2
		  and actual_at between $3 and $4
		order by actual_at desc, id desc
	`
	var rows []accrualDAO
	if err := pgxscan.Select(ctx, r.db.Pool(ctx), &rows, sql, targetUserID, targetRoleID, from, to); err != nil {
		return nil, err
	}

	result := make([]dto.Accrual, 0, len(rows))
	for _, row := range rows {
		result = append(result, dto.Accrual{
			ID:           row.ID,
			TargetUserID: row.TargetUserID,
			TargetRoleID: dto.Role(row.TargetRoleID),
			ActualTypeID: dto.AccrualActualType(row.ActualTypeID),
			LessonID:     derefInt64(row.LessonID),
			Amount:       row.Amount,
			Comment:      row.Comment,
			CreatedByID:  derefInt64(row.CreatedByID),
			ActualAt:     row.ActualAt,
			IsPaid:       row.IsPaid,
		})
	}

	return result, nil
}

func (r *Repository) CanManageTutor(ctx context.Context, tutorID, requesterUserID, requesterAdminID int64, isAssistant bool) (bool, error) {
	if !isAssistant {
		sql := `select exists(select 1 from tutors where id = $1 and admin_id = $2)`
		var exists bool
		err := pgxscan.Get(ctx, r.db.Pool(ctx), &exists, sql, tutorID, requesterAdminID)
		return exists, err
	}

	sql := `
		select exists(
			select 1
			from tutors t
			join assistant_tgs at on at.user_id = $2
			where t.id = $1
			  and t.tg_admin_username_id is not null
			  and t.tg_admin_username_id = any(at.available_tg_ids)
		)
	`
	var exists bool
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &exists, sql, tutorID, requesterUserID)
	return exists, err
}

func (r *Repository) CanManageAssistant(ctx context.Context, assistantID, requesterUserID, requesterAdminID int64, isAssistant bool) (bool, error) {
	if !isAssistant {
		sql := `select exists(select 1 from users where id = $1 and admin_id = $2 and role_id = $3)`
		var exists bool
		err := pgxscan.Get(ctx, r.db.Pool(ctx), &exists, sql, assistantID, requesterAdminID, dto.AssistantRole)
		return exists, err
	}

	sql := `
		select exists(
			select 1
			from assistant_tgs at
			where at.user_id = $1
			  and $2 = any(at.can_penalize_assistant_ids)
		)
	`
	var exists bool
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &exists, sql, requesterUserID, assistantID)
	return exists, err
}

func Summarize(accruals []dto.Accrual) dto.AccrualSummary {
	summary := dto.AccrualSummary{}
	for _, accrual := range accruals {
		switch accrual.ActualTypeID {
		case dto.AccrualActualTypeLesson:
			summary.Lessons = summary.Lessons.Add(accrual.Amount)
		case dto.AccrualActualTypePenalty:
			summary.Penalties = summary.Penalties.Add(accrual.Amount.Abs())
		case dto.AccrualActualTypeBonus:
			summary.Bonuses = summary.Bonuses.Add(accrual.Amount)
		}
		summary.Payable = summary.Payable.Add(accrual.Amount)
		if !accrual.IsPaid {
			summary.Unpaid = summary.Unpaid.Add(accrual.Amount)
		}
	}
	return summary
}

func derefInt64(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}
