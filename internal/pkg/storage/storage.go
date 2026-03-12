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
	Upload(ctx context.Context, bucket, key, contentType string, body io.Reader) error
	Download(ctx context.Context, bucket, key string) (*DownloadedObject, error)
	Delete(ctx context.Context, bucket, key string) error
}
