package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorage struct {
	client     *minio.Client
	bucketName string
}

// NewMinioStorage инициализирует клиент и проверяет наличие бакета
func NewMinioStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*FileStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	storage := &FileStorage{
		client:     client,
		bucketName: bucket,
	}

	// Создаем бакет при инициализации, если его нет
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return storage, nil
}

// UploadImage загружает файл из io.Reader (универсальный способ для веба)
func (s *FileStorage) UploadImage(ctx context.Context, fileName string, reader io.Reader, fileSize int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucketName, fileName, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// GetFileURL генерирует временную ссылку для просмотра файла
func (s *FileStorage) GetFileURL(ctx context.Context, fileName string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, fileName, expires, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// DeleteFile удаляет объект из хранилища
func (s *FileStorage) DeleteFile(ctx context.Context, fileName string) error {
	return s.client.RemoveObject(ctx, s.bucketName, fileName, minio.RemoveObjectOptions{})
}
