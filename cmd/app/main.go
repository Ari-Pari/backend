package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ari-Pari/backend/internal/api"
	generated "github.com/Ari-Pari/backend/internal/api/generated"
	"github.com/Ari-Pari/backend/internal/clients/dbstorage"
	"github.com/Ari-Pari/backend/internal/clients/filestorage"
	"github.com/Ari-Pari/backend/internal/config"
	"github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	dbPool := setupDbStorage(ctx, cfg)
	defer dbPool.Close()

	queries := db.New(dbPool.Pool)

	minioStore := setupMinioStorage(cfg)

	logger := log.New(os.Stdout, "API: ", log.LstdFlags|log.Lshortfile)

	server := api.NewServer(logger, queries, minioStore)

	router := setupRouter(server, logger)

	startServer(router, ":8080", logger)
}

func setupRouter(apiHandler *api.Server, logger *log.Logger) *chi.Mux {
	r := chi.NewRouter()

	// Базовые middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(5))

	r.Route("/api/v1", func(r chi.Router) {
		// Монтируем сгенерированный хендлер
		r.Mount("/", generated.Handler(apiHandler))
	})

	return r
}

func startServer(handler http.Handler, addr string, logger *log.Logger) {
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorLog:     logger,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Printf("Dance API starting on %s", addr)
		logger.Printf("Health check: http://localhost%s/health", addr)
		logger.Printf("Ready check: http://localhost%s/ready", addr)
		logger.Printf("API endpoints available under /api/v1")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-stop
	logger.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Printf("Server shutdown error: %v", err)
	}

	logger.Println("Server stopped")
}

func setupDbStorage(ctx context.Context, cfg *config.Config) *dbstorage.Storage {
	storage, err := dbstorage.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Successfully connected to the database from akob main!")
	return storage
}

func setupMinioStorage(cfg *config.Config) filestorage.FileStorage {
	fileStore, err := filestorage.NewMinioStorage(
		cfg.Minio.Endpoint,
		cfg.Minio.AccessKey,
		cfg.Minio.SecretKey,
		cfg.Minio.Bucket,
		false,
	)
	if err != nil {
		log.Fatalf("Failed to initialize file storage: %v", err)
	}
	log.Println("Successfully connected to MinIO!")
	return fileStore
}
