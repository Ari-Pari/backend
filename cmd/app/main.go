package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/Ari-Pari/backend/internal/clients/dbstorage"
	"github.com/Ari-Pari/backend/internal/clients/filestorage" // Твой новый импорт
	"github.com/Ari-Pari/backend/internal/config"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	// Инициализация БД
	storage, err := dbstorage.New(ctx, cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()
	log.Println("✅ Successfully connected to the database!")

	// Инициализация MinIO
	fileStore, err := filestorage.NewMinioStorage(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		false,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize file storage: %v", err)
	} else {
		log.Println("✅ Successfully connected to MinIO!")
		_ = fileStore
	}

	select {}
}
