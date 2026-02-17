package main

import (
	"context"
	"log"

	"github.com/Ari-Pari/backend/internal/clients/dbstorage"
	"github.com/Ari-Pari/backend/internal/config"
	db "github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/Ari-Pari/backend/internal/parser"
	"github.com/Ari-Pari/backend/internal/services/autoUploadDataService"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()
	cfg, err := config.Load()

	setupDbStorage(ctx, cfg)

	conn, err := pgxpool.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	myParser := parser.NewJSONParser()

	service := autoUploadDataService.NewAutoUploadDataService(db.New(conn))

	states, err := myParser.ParseStatesFile("static/autouploaddata/states.json")

	if err != nil {
		log.Fatal("Failed to parse states:", err)
	}

	err = service.ClearAllTables(ctx)

	if err != nil {
		log.Fatal("Failed to clear all tables:", err)
	}

	regions := parser.ToDomainRegions(states)

	err = service.CreateRegions(ctx, regions)
	if err != nil {
		log.Fatal("Failed to create regions:", err)
	}

	dances, err := myParser.ParseDancesFile("static/autouploaddata/dances.json")

	if err != nil {
		log.Fatal("Failed to parse dances:", err)
	}

	domainDances := parser.ToDomainDances(dances)

	err = service.CreateDances(ctx, domainDances)
	if err != nil {
		log.Fatal("Failed to create dances:", err)
	}

	musics, err := myParser.ParseMusicsFile("static/autouploaddata/musics.json")

	if err != nil {
		log.Fatal("Failed to parse musics:", err)
	}

	domainSongs := parser.ToDomainSongs(musics)

	err = service.CreateSongs(ctx, domainSongs)
	if err != nil {
		log.Fatal("Failed to create songs:", err)
	}

	videos, err := myParser.ParseVideosFile("static/autouploaddata/videos.json")

	if err != nil {
		log.Fatal("Failed to parse videos:", err)
	}

	domainVideos := parser.ToDomainVideos(videos)

	err = service.CreateVideos(ctx, domainVideos)
	if err != nil {
		log.Fatal("Failed to create videos:", err)
	}
}

func setupDbStorage(ctx context.Context, cfg *config.Config) {
	storage, err := dbstorage.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()
	log.Println("Successfully connected to the database!")
}
