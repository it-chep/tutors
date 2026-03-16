package storage

import (
	"context"
	"io"
)

type DownloadedObject struct {
	Body        io.ReadCloser
	ContentType string
}

type Storage interface {
	UploadContract(ctx context.Context, adminID, tutorID int64, fileName, contentType string, body io.Reader) (string, error)
	DownloadContract(ctx context.Context, key string) (*DownloadedObject, error)
	DeleteContract(ctx context.Context, key string) error
	UploadReceipt(ctx context.Context, adminID, tutorID int64, fileName, contentType string, body io.Reader) (string, error)
	DownloadReceipt(ctx context.Context, key string) (*DownloadedObject, error)
}
