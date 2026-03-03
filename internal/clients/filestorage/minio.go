package filestorage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// FileStorage описывает поведение нашего хранилища
type FileStorage interface {
	UploadFile(ctx context.Context, originalName string, reader io.Reader, fileSize int64, contentType string) (string, error)
	GetFileURL(fileKey string) (string, error)
	DeleteFile(ctx context.Context, fileKey string) error
	GetOriginalName(ctx context.Context, fileKey string) (string, error)
}

// minioStorage — внутренняя реализация интерфейса для MinIO
type minioStorage struct {
	client, publicClient *minio.Client
	publicURL            string
	bucketName           string
}

// NewMinioStorage возвращает интерфейс FileStorage
func NewMinioStorage(ctx context.Context, endpoint, serverURL, accessKey, secretKey, bucket string, useSSL bool) (FileStorage, error) {
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return &minioStorage{}, err
	}

	if !exists {
		// Создаем бакет
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return &minioStorage{}, err
		}
		log.Printf("Bucket '%s' created successfully", bucket)
	} else {
		log.Printf("Bucket '%s' already exists", bucket)
	}
	return &minioStorage{
		client:     client,
		bucketName: bucket,
		publicURL:  serverURL,
	}, nil
}

func (s *minioStorage) UploadFile(ctx context.Context, originalName string, reader io.Reader, fileSize int64, contentType string) (string, error) {
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

func (s *minioStorage) GetFileURL(fileKey string) (string, error) {
	return url.JoinPath(s.publicURL, s.bucketName, fileKey)
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
