package storage

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var ErrNotFound = errors.New("object not found")

type Presigner interface {
	PresignedGetURL(ctx context.Context, key, filename string) (string, error)
}

type S3Store struct {
	client  *s3.Client
	bucket  string
	presign *s3.PresignClient
	expires time.Duration
}

func NewS3Store(client *s3.Client, bucket string, expires time.Duration) *S3Store {
	return &S3Store{
		client:  client,
		bucket:  bucket,
		presign: s3.NewPresignClient(client),
		expires: expires,
	}
}

func (s *S3Store) PresignedGetURL(ctx context.Context, key, filename string) (string, error) {
	if _, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}); err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	req, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     aws.String(s.bucket),
		Key:                        aws.String(key),
		ResponseContentType:        aws.String(mimeFor(filename)),
		ResponseContentDisposition: aws.String(`attachment; filename="` + filepath.Base(filename) + `"`),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = s.expires
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func mimeFor(filename string) string {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".txt":
		return "text/plain"
	case ".zip":
		return "application/zip"
	case ".rar":
		return "application/x-rar-compressed"
	default:
		return "application/octet-stream"
	}
}
