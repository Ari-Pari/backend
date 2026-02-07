package filestorage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// FileStorage описывает поведение нашего хранилища
type FileStorage interface {
	UploadImage(ctx context.Context, originalName string, reader io.Reader, fileSize int64, contentType string) (string, error)
	GetFileURL(ctx context.Context, fileKey string, expires time.Duration) (string, error)
	DeleteFile(ctx context.Context, fileKey string) error
	GetOriginalName(ctx context.Context, fileKey string) (string, error)
}

// minioStorage — внутренняя реализация интерфейса для MinIO
type minioStorage struct {
	client     *minio.Client
	bucketName string
}

// NewMinioStorage возвращает интерфейс FileStorage
func NewMinioStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (FileStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &minioStorage{
		client:     client,
		bucketName: bucket,
	}, nil
}

func (s *minioStorage) UploadImage(ctx context.Context, originalName string, reader io.Reader, fileSize int64, contentType string) (string, error) {
	ext := filepath.Ext(originalName)
	fileKey := uuid.New().String() + ext

	opts := minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"Original-Name": originalName,
		},
	}

	_, err := s.client.PutObject(ctx, s.bucketName, fileKey, reader, fileSize, opts)
	if err != nil {
		return "", err
	}

	return fileKey, nil
}

func (s *minioStorage) GetFileURL(ctx context.Context, fileKey string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)

	origName, _ := s.GetOriginalName(ctx, fileKey)
	if origName != "" {
		reqParams.Set("response-content-disposition", fmt.Sprintf("inline; filename=\"%s\"", origName))
	} else {
		reqParams.Set("response-content-disposition", "inline")
	}

	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, fileKey, expires, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (s *minioStorage) GetOriginalName(ctx context.Context, fileKey string) (string, error) {
	info, err := s.client.StatObject(ctx, s.bucketName, fileKey, minio.StatObjectOptions{})
	if err != nil {
		return "", err
	}
	return info.UserMetadata["Original-Name"], nil
}

func (s *minioStorage) DeleteFile(ctx context.Context, fileKey string) error {
	return s.client.RemoveObject(ctx, s.bucketName, fileKey, minio.RemoveObjectOptions{})
}
