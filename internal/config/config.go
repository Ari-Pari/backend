package config

import (
	"fmt"
	"os"
)

type PostgresConfig struct {
	DSN string
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type Config struct {
	Postgres PostgresConfig
	Minio    MinioConfig
}

func Load() (*Config, error) {
	pgConfig, err := loadPostgresConfig()
	if err != nil {
		return nil, err
	}

	minioConfig, err := loadMinioConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Postgres: *pgConfig,
		Minio:    *minioConfig,
	}, nil
}

func loadPostgresConfig() (*PostgresConfig, error) {
	user, err := getEnv("POSTGRES_USER")
	if err != nil {
		return nil, err
	}

	password, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, err
	}

	host, err := getEnv("POSTGRES_HOST")
	if err != nil {
		return nil, err
	}

	port, err := getEnv("POSTGRES_PORT")
	if err != nil {
		return nil, err
	}

	dbName, err := getEnv("POSTGRES_DB")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName,
	)

	return &PostgresConfig{
		DSN: dsn,
	}, nil
}

func loadMinioConfig() (*MinioConfig, error) {
	endpoint, err := getEnv("MINIO_ENDPOINT")
	if err != nil {
		return nil, err
	}

	accessKey, err := getEnv("MINIO_ACCESS_KEY")
	if err != nil {
		return nil, err
	}

	secretKey, err := getEnv("MINIO_SECRET_KEY")
	if err != nil {
		return nil, err
	}

	bucket, err := getEnv("MINIO_BUCKET")
	if err != nil {
		return nil, err
	}

	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	return &MinioConfig{
		Endpoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		UseSSL:    useSSL,
	}, nil
}

func getEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return val, nil
}
