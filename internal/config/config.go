package config

import (
	"fmt"
	"os"
)

type Config struct {
	DSN string
}

func Load() (*Config, error) {
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

	return &Config{
		DSN: dsn,
	}, nil
}