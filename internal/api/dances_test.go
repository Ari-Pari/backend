package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	api "github.com/Ari-Pari/backend/internal/api/generated"
	"github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testDBPool *pgxpool.Pool

type mockStorage struct{}

func (m *mockStorage) GetFileURL(key string) (string, error) {
	return "http://minio/" + key, nil
}
func (m *mockStorage) UploadFile(ctx context.Context, originalName string, reader io.Reader, fileSize int64, contentType string) (string, error) {
	return "", nil
}
func (m *mockStorage) DeleteFile(context.Context, string) error                { return nil }
func (m *mockStorage) GetOriginalName(context.Context, string) (string, error) { return "", nil }

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := testpostgres.Run(ctx,
		"postgres:17-alpine",
		testpostgres.WithDatabase("aripari_test"),
		testpostgres.WithUsername("user"),
		testpostgres.WithPassword("password"),
		testpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	mig, err := migrate.New("file://../../migrations", connStr)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrate up: %v", err)
	}

	testDBPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer testDBPool.Close()

	os.Exit(m.Run())
}

func clearTables(t *testing.T) {
	ctx := context.Background()
	_, err := testDBPool.Exec(ctx, `
		TRUNCATE TABLE dance_song, songs, dance_region, videos, dance_videos, dances, regions, translations RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}

func TestGetDancesId_Integration(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	var translationID int64
	err := testDBPool.QueryRow(ctx, "INSERT INTO translations (eng_name, ru_name) VALUES ('Berd', 'Берд') RETURNING id").Scan(&translationID)
	require.NoError(t, err)

	_, err = testDBPool.Exec(ctx, `
		INSERT INTO dances (id, translation_id, name, complexity, photo_key, gender, paces, genres, handshakes) 
		VALUES (1, $1, 'Berd_def', 3, 'photo.jpg', 'male', '{1,2}', '{"WAR", "FESTIVE"}', '{"SHOULDER"}')
	`, translationID)
	require.NoError(t, err)

	var regTransID int64
	err = testDBPool.QueryRow(ctx, "INSERT INTO translations (eng_name, ru_name) VALUES ('Shirak', 'Ширак') RETURNING id").Scan(&regTransID)
	require.NoError(t, err)

	_, err = testDBPool.Exec(ctx, "INSERT INTO regions (id, translation_id, name) VALUES (10, $1, 'Shirak_def')", regTransID)
	require.NoError(t, err)
	_, err = testDBPool.Exec(ctx, "INSERT INTO dance_region (dance_id, region_id) VALUES (1, 10)")
	require.NoError(t, err)

	var videoTransID int64
	err = testDBPool.QueryRow(ctx, "INSERT INTO translations (eng_name, ru_name) VALUES ('Berd Video', 'Видео Берд') RETURNING id").Scan(&videoTransID)
	require.NoError(t, err)

	_, err = testDBPool.Exec(ctx, "INSERT INTO videos (id, translation_id, name, link, type) VALUES (100, $1, 'Video_def', 'http://yt', 'source')", videoTransID)
	require.NoError(t, err)
	_, err = testDBPool.Exec(ctx, "INSERT INTO dance_videos (dance_id, video_id) VALUES (1, 100)")
	require.NoError(t, err)

	var songTransID int64
	err = testDBPool.QueryRow(ctx, "INSERT INTO translations (eng_name, ru_name) VALUES ('Berd Song', 'Песня Берд') RETURNING id").Scan(&songTransID)
	require.NoError(t, err)

	_, err = testDBPool.Exec(ctx, "INSERT INTO songs (id, translation_id, file_key, name) VALUES (50, $1, 'song.mp3', 'Berd_Song_def')", songTransID)
	require.NoError(t, err)
	_, err = testDBPool.Exec(ctx, "INSERT INTO dance_song (dance_id, song_id) VALUES (1, 50)")
	require.NoError(t, err)

	var artistTransID int64
	testDBPool.QueryRow(ctx, "INSERT INTO translations (eng_name, ru_name) VALUES ('Ensemble A', 'Ансамбль А') RETURNING id").Scan(&artistTransID)
	
	testDBPool.Exec(ctx, "INSERT INTO artists (id, translation_id, name, link) VALUES (200, $1, 'Ens_def', 'http://ens.com')", artistTransID)
	testDBPool.Exec(ctx, "INSERT INTO song_artist (song_id, artist_id) VALUES (50, 200)")

	// запуск теста
	queries := db.New(testDBPool)
	logger := log.New(io.Discard, "", 0)
	srv := NewServer(logger, queries, &mockStorage{})

	// без языка
	t.Run("Success 200 - Default Lang", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/dances/1", nil)
		w := httptest.NewRecorder()

		srv.GetDancesId(w, req, 1, api.GetDancesIdParams{})

		assert.Equal(t, http.StatusOK, w.Code)

		var response api.DanceFullResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Berd_def", response.Name)
		assert.Equal(t, "Shirak_def", response.Regions[0].Name)
	})

	// русский язык
	t.Run("Success 200 - Russian Lang", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/dances/1?lang=ru", nil)
		w := httptest.NewRecorder()

		lang := "ru"
		srv.GetDancesId(w, req, 1, api.GetDancesIdParams{Lang: &lang})

		assert.Equal(t, http.StatusOK, w.Code)

		var response api.DanceFullResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Проверка танца
		assert.Equal(t, "Берд", response.Name)

		// Проверка региона
		require.Len(t, response.Regions, 1)
		assert.Equal(t, "Ширак", response.Regions[0].Name)

		// Проверка видео
		require.NotNil(t, response.SourceVideos)
		require.Len(t, *response.SourceVideos, 1)
		assert.Equal(t, "Видео Берд", (*response.SourceVideos)[0].Name)

		// Проверка песни
		require.Len(t, response.Songs, 1)
		assert.Equal(t, "Песня Берд", response.Songs[0].Name)

		require.Len(t, response.Songs[0].Ensembles, 1)
    	assert.Equal(t, "Ансамбль А", response.Songs[0].Ensembles[0].Name)
	})

	// ТЕСТ 3: НЕ НАЙДЕНО
	t.Run("Not Found 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/dances/999", nil)
		w := httptest.NewRecorder()

		srv.GetDancesId(w, req, 999, api.GetDancesIdParams{})

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetRegions_Integration(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	// Регион 1: Ширак
	var transID1 int64
	err := testDBPool.QueryRow(ctx, "INSERT INTO translations (eng_name, ru_name) VALUES ('Shirak', 'Ширак') RETURNING id").Scan(&transID1)
	require.NoError(t, err)
	_, err = testDBPool.Exec(ctx, "INSERT INTO regions (id, translation_id, name) VALUES (1, $1, 'Shirak_default')", transID1)
	require.NoError(t, err)

	// Регион 2: Лори
	var transID2 int64
	err = testDBPool.QueryRow(ctx, "INSERT INTO translations (eng_name, ru_name) VALUES ('Lori', 'Лори') RETURNING id").Scan(&transID2)
	require.NoError(t, err)
	_, err = testDBPool.Exec(ctx, "INSERT INTO regions (id, translation_id, name) VALUES (2, $1, 'Lori_default')", transID2)
	require.NoError(t, err)

	// сервер
	queries := db.New(testDBPool)
	logger := log.New(io.Discard, "", 0)
	srv := NewServer(logger, queries, &mockStorage{})

	t.Run("Success 200 - Russian Lang", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/regions?lang=ru", nil)
		w := httptest.NewRecorder()

		lang := "ru"
		srv.GetRegions(w, req, api.GetRegionsParams{Lang: &lang})

		assert.Equal(t, http.StatusOK, w.Code)

		var response api.RegionListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Len(t, response, 2)
		assert.Equal(t, "Ширак", response[0].Name)
		assert.Equal(t, "Лори", response[1].Name)
	})

	t.Run("Success 200 - Fallback Lang", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/regions?lang=fr", nil)
		w := httptest.NewRecorder()

		lang := "fr"
		srv.GetRegions(w, req, api.GetRegionsParams{Lang: &lang})

		assert.Equal(t, http.StatusOK, w.Code)

		var response api.RegionListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Len(t, response, 2)
		assert.Equal(t, "Shirak_default", response[0].Name)
		assert.Equal(t, "Lori_default", response[1].Name)
	})
}
