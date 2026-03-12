package contract

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/action/contract/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/storage"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Action struct {
	dal     *dal.Repository
	storage storage.Storage
	bucket  string
}

type UploadRequest struct {
	FileName    string
	ContentType string
	Body        io.Reader
}

type Contract struct {
	TutorID     int64
	TutorName   string
	FileKey     string
	FileName    string
	ContentType string
	CreatedAt   time.Time
}

type DownloadedContract struct {
	Contract
	Body io.ReadCloser
}

func New(pool *pgxpool.Pool, objectStorage storage.Storage, bucket string) *Action {
	return &Action{
		dal:     dal.NewRepository(pool),
		storage: objectStorage,
		bucket:  bucket,
	}
}

func (a *Action) Upload(ctx context.Context, tutorID int64, req UploadRequest) (Contract, error) {
	if err := a.ensureWriteAccess(ctx, tutorID); err != nil {
		return Contract{}, err
	}
	if req.FileName == "" {
		return Contract{}, errors.New("file name is required")
	}
	if req.Body == nil {
		return Contract{}, errors.New("file body is required")
	}

	contentType := req.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	key := buildObjectKey("contracts", tutorID, req.FileName)
	if err := a.storage.Upload(ctx, a.bucket, key, contentType, req.Body); err != nil {
		return Contract{}, err
	}

	uploadedByID := userCtx.UserIDFromContext(ctx)
	if err := a.dal.Upsert(ctx, tutorID, key, req.FileName, contentType, uploadedByID); err != nil {
		return Contract{}, err
	}

	contract, err := a.dal.Get(ctx, tutorID)
	if err != nil {
		return Contract{}, err
	}

	return toDomain(contract), nil
}

func (a *Action) Get(ctx context.Context, tutorID int64) (DownloadedContract, error) {
	if err := a.ensureReadAccess(ctx, tutorID); err != nil {
		return DownloadedContract{}, err
	}

	contract, err := a.dal.Get(ctx, tutorID)
	if err != nil {
		return DownloadedContract{}, err
	}

	downloaded, err := a.storage.Download(ctx, a.bucket, contract.FileKey)
	if err != nil {
		return DownloadedContract{}, err
	}

	domain := toDomain(contract)
	if downloaded.ContentType != "" {
		domain.ContentType = downloaded.ContentType
	}

	return DownloadedContract{
		Contract: domain,
		Body:     downloaded.Body,
	}, nil
}

func (a *Action) Delete(ctx context.Context, tutorID int64) error {
	if err := a.ensureWriteAccess(ctx, tutorID); err != nil {
		return err
	}

	contract, err := a.dal.Get(ctx, tutorID)
	if err != nil {
		return err
	}

	if err = a.dal.Delete(ctx, tutorID); err != nil {
		return err
	}

	if err = a.storage.Delete(ctx, a.bucket, contract.FileKey); err != nil {
		return err
	}

	return nil
}

func (a *Action) ListVisible(ctx context.Context) ([]Contract, error) {
	contracts, err := a.listVisibleRecords(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Contract, 0, len(contracts))
	for _, contract := range contracts {
		result = append(result, toDomain(contract))
	}

	return result, nil
}

func (a *Action) DownloadByKey(ctx context.Context, key string) (*storage.DownloadedObject, error) {
	return a.storage.Download(ctx, a.bucket, key)
}

func (a *Action) ensureReadAccess(ctx context.Context, tutorID int64) error {
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
		canViewContracts, err := a.dal.AssistantCanViewContracts(ctx, userCtx.UserIDFromContext(ctx))
		if err != nil {
			return err
		}
		if !canViewContracts {
			return errors.New("access denied")
		}

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

func (a *Action) ensureWriteAccess(ctx context.Context, tutorID int64) error {
	if dto.IsSuperAdminRole(ctx) {
		return nil
	}
	if !dto.IsAdminRole(ctx) {
		return errors.New("access denied")
	}

	ok, err := a.dal.CanAdminManageTutor(ctx, tutorID, userCtx.AdminIDFromContext(ctx))
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("access denied")
	}

	return nil
}

func (a *Action) listVisibleRecords(ctx context.Context) ([]dal.Contract, error) {
	switch {
	case dto.IsSuperAdminRole(ctx):
		return a.dal.ListAll(ctx)
	case dto.IsAdminRole(ctx):
		return a.dal.ListByAdmin(ctx, userCtx.AdminIDFromContext(ctx))
	case dto.IsAssistantRole(ctx):
		canViewContracts, err := a.dal.AssistantCanViewContracts(ctx, userCtx.UserIDFromContext(ctx))
		if err != nil {
			return nil, err
		}
		if !canViewContracts {
			return nil, errors.New("access denied")
		}
		return a.dal.ListByAssistant(ctx, userCtx.UserIDFromContext(ctx), userCtx.AdminIDFromContext(ctx))
	default:
		return nil, errors.New("access denied")
	}
}

func toDomain(contract dal.Contract) Contract {
	return Contract{
		TutorID:     contract.TutorID,
		TutorName:   contract.TutorName,
		FileKey:     contract.FileKey,
		FileName:    contract.FileName,
		ContentType: contract.ContentType,
		CreatedAt:   contract.CreatedAt,
	}
}

func buildObjectKey(prefix string, targetID int64, fileName string) string {
	safeName := strings.TrimSpace(filepath.Base(fileName))
	if safeName == "." || safeName == "/" || safeName == "" {
		safeName = "file"
	}

	replacer := strings.NewReplacer(" ", "_", "/", "_", "\\", "_", ":", "_")
	safeName = replacer.Replace(safeName)

	return fmt.Sprintf("%s/%d/%d_%s", prefix, targetID, time.Now().UTC().UnixNano(), safeName)
}

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
