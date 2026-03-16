package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/it-chep/tutors.git/internal/config"
)

type s3Storage struct {
	client          *s3.Client
	bucketContracts string
	bucketReceipts  string
}

func NewS3(cfg config.S3Config) (Storage, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
		if !strings.EqualFold(service, s3.ServiceID) {
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		}

		return aws.Endpoint{
			URL:               cfg.Endpoint,
			HostnameImmutable: true,
			SigningRegion:     cfg.Region,
		}, nil
	})

	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		awsConfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("load s3 config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &s3Storage{
		client:          client,
		bucketContracts: cfg.ContractsBucket,
		bucketReceipts:  cfg.ReceiptsBucket,
	}, nil
}

func (s *s3Storage) UploadContract(ctx context.Context, adminID, tutorID int64, fileName, contentType string, body io.Reader) (string, error) {
	key := s.buildObjectKey(adminID, tutorID, fileName)
	if err := s.upload(ctx, s.bucketContracts, key, contentType, body); err != nil {
		return "", err
	}

	return key, nil
}

func (s *s3Storage) DownloadContract(ctx context.Context, key string) (*DownloadedObject, error) {
	return s.download(ctx, s.bucketContracts, key)
}

func (s *s3Storage) DeleteContract(ctx context.Context, key string) error {
	return s.delete(ctx, s.bucketContracts, key)
}

func (s *s3Storage) UploadReceipt(ctx context.Context, adminID, tutorID int64, fileName, contentType string, body io.Reader) (string, error) {
	key := s.buildObjectKey(adminID, tutorID, fileName)
	if err := s.upload(ctx, s.bucketReceipts, key, contentType, body); err != nil {
		return "", err
	}

	return key, nil
}

func (s *s3Storage) DownloadReceipt(ctx context.Context, key string) (*DownloadedObject, error) {
	return s.download(ctx, s.bucketReceipts, key)
}

func (s *s3Storage) upload(ctx context.Context, bucket, key, contentType string, body io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
		Metadata: map[string]string{
			"uploaded_at": time.Now().Format(time.RFC3339),
			"origin_name": "",
		},
	})

	return err
}

func (s *s3Storage) download(ctx context.Context, bucket, key string) (*DownloadedObject, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return &DownloadedObject{
		Body:        out.Body,
		ContentType: aws.ToString(out.ContentType),
	}, nil
}

func (s *s3Storage) delete(ctx context.Context, bucket, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *s3Storage) buildObjectKey(adminID, tutorID int64, fileName string) string {
	safeName := strings.TrimSpace(filepath.Base(fileName))
	if safeName == "." || safeName == "/" || safeName == "" {
		safeName = "file"
	}

	replacer := strings.NewReplacer(" ", "_", "/", "_", "\\", "_", ":", "_")
	safeName = replacer.Replace(safeName)

	return fmt.Sprintf("admin_%d/tutor_%d/%s", adminID, tutorID, safeName)
}
