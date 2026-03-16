package payout

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/it-chep/tutors.git/internal/module/admin/action/payout/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/storage"
	"github.com/it-chep/tutors.git/internal/pkg/transaction"
	"github.com/it-chep/tutors.git/internal/pkg/transaction/wrapper"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Action struct {
	dal     *dal.Repository
	storage storage.Storage
}

type CreateRequest struct {
	TutorID int64
	Amount  int64
	Comment string
}

type ReceiptUploadRequest struct {
	FileName    string
	ContentType string
	Body        io.Reader
}

type Payout struct {
	ID                 uuid.UUID
	TutorID            int64
	TutorName          string
	Amount             decimal.Decimal
	Comment            string
	CreatedByID        int64
	CreatedAt          time.Time
	ReceiptFileName    string
	ReceiptContentType string
	ReceiptUploadedAt  *time.Time
}

type DownloadedReceipt struct {
	Payout
	Body io.ReadCloser
}

func New(pool *pgxpool.Pool, objectStorage storage.Storage) *Action {
	return &Action{
		dal:     dal.NewRepository(wrapper.NewDatabase(pool)),
		storage: objectStorage,
	}
}

func (a *Action) Create(ctx context.Context, req CreateRequest) (uuid.UUID, error) {
	if err := a.ensureManageTutor(ctx, req.TutorID); err != nil {
		return uuid.Nil, err
	}
	if req.Amount <= 0 {
		return uuid.Nil, errors.New("amount must be greater than zero")
	}

	payoutID := uuid.New()
	createdByID := userCtx.UserIDFromContext(ctx)

	err := transaction.Exec(ctx, func(ctx context.Context) error {
		if err := a.dal.Create(ctx, dal.CreatePayoutRequest{
			ID:           payoutID,
			TargetUserID: req.TutorID,
			TargetRoleID: int64(dto.TutorRole),
			CreatedByID:  createdByID,
			Amount:       decimal.NewFromInt(req.Amount),
			Comment:      req.Comment,
		}); err != nil {
			return err
		}

		updatedCount, err := a.dal.MarkTutorAccrualsPaid(ctx, req.TutorID, payoutID)
		if err != nil {
			return err
		}
		if updatedCount == 0 {
			return errors.New("no unpaid accruals")
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return payoutID, nil
}

func (a *Action) ListReceipts(ctx context.Context, tutorID int64, from, to time.Time) ([]Payout, error) {
	if err := a.ensureVisibleTutor(ctx, tutorID); err != nil {
		return nil, err
	}

	items, err := a.dal.ListTutorPayouts(ctx, tutorID, from, to)
	if err != nil {
		return nil, err
	}

	return toDomainList(items), nil
}

func (a *Action) ListVisibleReceipts(ctx context.Context, from, to time.Time) ([]Payout, error) {
	var items []dal.Payout
	var err error

	switch {
	case dto.IsSuperAdminRole(ctx):
		items, err = a.dal.ListAllPayouts(ctx, from, to)
	case dto.IsAdminRole(ctx):
		items, err = a.dal.ListAdminPayouts(ctx, userCtx.AdminIDFromContext(ctx), from, to)
	case dto.IsAssistantRole(ctx):
		items, err = a.dal.ListAssistantPayouts(ctx, userCtx.UserIDFromContext(ctx), userCtx.AdminIDFromContext(ctx), from, to)
	default:
		return nil, errors.New("access denied")
	}
	if err != nil {
		return nil, err
	}

	return toDomainList(items), nil
}

func (a *Action) SaveReceipt(ctx context.Context, req ReceiptUploadRequest) (uuid.UUID, error) {
	if !dto.IsTutorRole(ctx) {
		return uuid.Nil, errors.New("access denied")
	}
	if req.FileName == "" {
		return uuid.Nil, errors.New("file name is required")
	}
	if req.Body == nil {
		return uuid.Nil, errors.New("file body is required")
	}

	tutorID := userCtx.GetTutorID(ctx)
	payout, err := a.dal.GetNearestPayoutWithoutReceipt(ctx, tutorID)
	if err != nil {
		return uuid.Nil, err
	}

	contentType := req.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	adminID := userCtx.AdminIDFromContext(ctx)

	key, err := a.storage.UploadReceipt(ctx, adminID, tutorID, req.FileName, contentType, req.Body)
	if err != nil {
		return uuid.Nil, err
	}

	if err = a.dal.AttachReceipt(ctx, payout.ID, key, req.FileName, contentType); err != nil {
		return uuid.Nil, err
	}

	return payout.ID, nil
}

func (a *Action) DownloadReceipt(ctx context.Context, payoutID uuid.UUID) (DownloadedReceipt, error) {
	payout, err := a.dal.GetPayout(ctx, payoutID)
	if err != nil {
		return DownloadedReceipt{}, err
	}
	if payout.ReceiptFileName == "" || payout.ReceiptFileKey == "" {
		return DownloadedReceipt{}, errors.New("receipt not found")
	}
	if err = a.ensureVisibleTutor(ctx, payout.TargetUserID); err != nil {
		return DownloadedReceipt{}, err
	}

	downloaded, err := a.storage.DownloadReceipt(ctx, payout.ReceiptFileKey)
	if err != nil {
		return DownloadedReceipt{}, err
	}

	result := toDomain(payout)
	if downloaded.ContentType != "" {
		result.ReceiptContentType = downloaded.ContentType
	}

	return DownloadedReceipt{
		Payout: result,
		Body:   downloaded.Body,
	}, nil
}

func (a *Action) IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func (a *Action) ensureManageTutor(ctx context.Context, tutorID int64) error {
	switch {
	case dto.IsSuperAdminRole(ctx):
		return nil
	case dto.IsAdminRole(ctx):
		ok, err := a.dal.CanAdminManageTutor(ctx, tutorID, userCtx.AdminIDFromContext(ctx))
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("access denied")
		}
		return nil
	case dto.IsAssistantRole(ctx):
		ok, err := a.dal.CanAssistantAccessTutor(ctx, tutorID, userCtx.UserIDFromContext(ctx), userCtx.AdminIDFromContext(ctx))
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("access denied")
		}
		return nil
	default:
		return errors.New("access denied")
	}
}

func (a *Action) ensureVisibleTutor(ctx context.Context, tutorID int64) error {
	return a.ensureManageTutor(ctx, tutorID)
}

func toDomainList(items []dal.Payout) []Payout {
	result := make([]Payout, 0, len(items))
	for _, item := range items {
		result = append(result, toDomain(item))
	}
	return result
}

func toDomain(item dal.Payout) Payout {
	return Payout{
		ID:                 item.ID,
		TutorID:            item.TargetUserID,
		TutorName:          item.TutorName,
		Amount:             item.Amount,
		Comment:            item.Comment,
		CreatedByID:        item.CreatedByID,
		CreatedAt:          item.CreatedAt,
		ReceiptFileName:    item.ReceiptFileName,
		ReceiptContentType: item.ReceiptContentType,
		ReceiptUploadedAt:  item.ReceiptUploadedAt,
	}
}
