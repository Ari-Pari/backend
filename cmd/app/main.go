package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/Ari-Pari/backend/internal/clients/dbstorage"
	"github.com/Ari-Pari/backend/internal/config"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	storage, err := dbstorage.New(ctx, cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()

	log.Println("Successfully connected to the database!")

	select {}
}
