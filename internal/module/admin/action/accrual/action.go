package accrual

import (
	"context"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/action/accrual/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/transaction"
	"github.com/it-chep/tutors.git/internal/pkg/transaction/wrapper"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Action struct {
	dal *dal.Repository
}

type CreateRequest struct {
	TargetUserID int64
	TargetRoleID dto.Role
	Type         dto.AccrualActualType
	Amount       int64
	Comment      string
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(wrapper.NewDatabase(pool)),
	}
}

func (a *Action) Create(ctx context.Context, req CreateRequest) error {
	createdByID := userCtx.UserIDFromContext(ctx)

	switch req.TargetRoleID {
	case dto.TutorRole:
		allowed, err := a.dal.CanManageTutor(ctx, req.TargetUserID, userCtx.UserIDFromContext(ctx), userCtx.AdminIDFromContext(ctx), dto.IsAssistantRole(ctx))
		if err != nil {
			return err
		}
		if !allowed {
			return errors.New("access denied")
		}
	case dto.AssistantRole:
		allowed, err := a.dal.CanManageAssistant(ctx, req.TargetUserID, userCtx.UserIDFromContext(ctx), userCtx.AdminIDFromContext(ctx), dto.IsAssistantRole(ctx))
		if err != nil {
			return err
		}
		if !allowed {
			return errors.New("access denied")
		}
	default:
		return errors.New("unsupported target role")
	}

	amount := decimal.NewFromInt(req.Amount)
	if req.Type == dto.AccrualActualTypePenalty && amount.GreaterThan(decimal.Zero) {
		amount = amount.Neg()
	}
	if req.Type == dto.AccrualActualTypeBonus && amount.LessThan(decimal.Zero) {
		amount = amount.Abs()
	}

	return a.dal.Create(ctx, dal.CreateAccrualRequest{
		TargetUserID: req.TargetUserID,
		TargetRoleID: int64(req.TargetRoleID),
		ActualTypeID: int64(req.Type),
		Amount:       amount,
		Comment:      req.Comment,
		CreatedByID:  createdByID,
		ActualAt:     time.Now().UTC(),
	})
}

func (a *Action) List(ctx context.Context, targetUserID int64, targetRoleID dto.Role, from, to time.Time) ([]dto.Accrual, dto.AccrualSummary, error) {
	switch targetRoleID {
	case dto.TutorRole:
		allowed, err := a.dal.CanManageTutor(ctx, targetUserID, userCtx.UserIDFromContext(ctx), userCtx.AdminIDFromContext(ctx), dto.IsAssistantRole(ctx))
		if err != nil {
			return nil, dto.AccrualSummary{}, err
		}
		if !allowed {
			return nil, dto.AccrualSummary{}, errors.New("access denied")
		}
	case dto.AssistantRole:
		allowed, err := a.dal.CanManageAssistant(ctx, targetUserID, userCtx.UserIDFromContext(ctx), userCtx.AdminIDFromContext(ctx), dto.IsAssistantRole(ctx))
		if err != nil {
			return nil, dto.AccrualSummary{}, err
		}
		if !allowed && userCtx.UserIDFromContext(ctx) != targetUserID {
			return nil, dto.AccrualSummary{}, errors.New("access denied")
		}
	default:
		return nil, dto.AccrualSummary{}, errors.New("unsupported target role")
	}

	accruals, err := a.dal.List(ctx, targetUserID, int64(targetRoleID), from, to)
	if err != nil {
		return nil, dto.AccrualSummary{}, err
	}
	summary := dal.Summarize(accruals)
	return accruals, summary, nil
}

func (a *Action) UpsertLessonAccrual(ctx context.Context, lessonID, tutorID int64, amount decimal.Decimal, actualAt time.Time) error {
	return a.dal.UpsertLessonAccrual(ctx, lessonID, tutorID, amount, actualAt)
}

func (a *Action) DeleteLessonAccrual(ctx context.Context, lessonID int64) error {
	return a.dal.DeleteLessonAccrual(ctx, lessonID)
}

func (a *Action) Exec(ctx context.Context, callback func(ctx context.Context) error) error {
	return transaction.Exec(ctx, callback)
}
