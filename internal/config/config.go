package config

import (
	"fmt"
	"os"
)

type Config struct {
	DSN            string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

func Load() (*Config, error) {
	// --- Данные для PostgreSQL ---
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	if user == "" || password == "" || host == "" || port == "" || dbName == "" {
		return nil, fmt.Errorf("one or more environment variables (POSTGRES_*) are missing")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName,
	)

	// --- Данные для MinIO ---
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	minioBucket := os.Getenv("MINIO_BUCKET")
	minioUseSSLStr := os.Getenv("MINIO_USE_SSL")

	if minioEndpoint == "" || minioAccessKey == "" || minioSecretKey == "" || minioBucket == "" {
		return nil, fmt.Errorf("one or more environment variables (MINIO_*) are missing")
	}

	// Конвертируем строку в bool (если в .env написано "true", будет true)
	useSSL := minioUseSSLStr == "true"

	return &Config{
		DSN:            dsn,
		MinioEndpoint:  minioEndpoint,
		MinioAccessKey: minioAccessKey,
		MinioSecretKey: minioSecretKey,
		MinioBucket:    minioBucket,
		MinioUseSSL:    useSSL,
	}, nil
}
