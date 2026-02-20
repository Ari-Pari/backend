package autoUploadDataService

import (
	"context"
	"os"
	"testing"

	"github.com/Ari-Pari/backend/internal/config"
	db "github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/Ari-Pari/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var pool *pgxpool.Pool
var querier db.Querier

func TestMain(m *testing.M) {
	err := godotenv.Load(".env.test")
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	// Загружаем конфиг
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Подключаемся к БД
	pool, err = pgxpool.New(context.Background(), cfg.Postgres.DSN)
	if err != nil {
		panic("Failed to connect to PostgreSQL: " + err.Error())
	}

	querier = db.New(pool)

	// Запускаем тесты
	code := m.Run()

	// Завершаем подключение
	pool.Close()

	os.Exit(code)
}

func resetDB(t *testing.T) {
	err := querier.TruncateAllTables(context.Background())
	require.NoError(t, err)
}

// Тесты остаются без изменений
func TestCreateRegions_Integration(t *testing.T) {
	resetDB(t)

	service := NewAutoUploadDataService(querier)

	regions := []domain.Region{
		{
			Id: 1,
			Name: domain.Translation{
				ArmName: "Շիրակ",
				EngName: "Shirak",
				RuName:  "Ширак",
			},
		},
		{
			Id: 2,
			Name: domain.Translation{
				ArmName: "Լոռի",
				EngName: "Lori",
				RuName:  "Лори",
			},
		},
	}

	err := service.CreateRegions(context.Background(), regions)
	require.NoError(t, err)

	dbRegions, err := querier.GetRegions(context.Background())
	require.NoError(t, err)
	assert.Len(t, dbRegions, len(regions))
}

func TestCreateDances_Integration(t *testing.T) {
	resetDB(t)

	service := NewAutoUploadDataService(querier)

	dances := []domain.DanceShort{
		{
			Id:      1,
			NameKey: "dance.shirak",
			Name: domain.Translation{
				ArmName: "Շիրակ",
				EngName: "Shirak",
				RuName:  "Ширак",
			},
			Complexity: &[]int32{3}[0],
			Genres:     []domain.Genre{domain.War, domain.Lyrical},
			Gender:     domain.Male,
			Paces:      []int32{1, 2, 3},
			HoldingTypes: []domain.HoldingType{
				domain.Free,
				domain.Palm,
			},
			RegionIds: []int64{1, 2},
		},
	}

	err := service.CreateDances(context.Background(), dances)
	require.NoError(t, err)

	dbDances, err := querier.GetDances(context.Background())
	require.NoError(t, err)
	assert.Len(t, dbDances, len(dances))
}

func TestCreateSongs_Integration(t *testing.T) {
	resetDB(t)

	service := NewAutoUploadDataService(querier)

	songs := []domain.SongShort{
		{
			Id:      1,
			NameKey: "song.shirak",
			Name: domain.Translation{
				ArmName: "Շիրակյան երգ",
				EngName: "Song of Shirak",
				RuName:  "Ширакская песня",
			},
			DanceIds: []int64{1, 2},
		},
	}

	err := service.CreateSongs(context.Background(), songs)
	require.NoError(t, err)

	dbSongs, err := querier.GetSongs(context.Background())
	require.NoError(t, err)
	assert.Len(t, dbSongs, len(songs))
}

func TestCreateVideos_Integration(t *testing.T) {
	resetDB(t)

	service := NewAutoUploadDataService(querier)

	videos := []domain.VideoShort{
		{
			Id:      nil,
			NameKey: "video.shirak",
			Name: domain.Translation{
				ArmName: "Շիրակ տեսանյութ",
				EngName: "Shirak video",
				RuName:  "Ширакское видео",
			},
			Link:     "https://example.com/video1",
			DanceIds: []int64{1, 2},
			Type:     domain.Lesson,
		},
	}

	err := service.CreateVideos(context.Background(), videos)
	require.NoError(t, err)

	dbVideos, err := querier.GetVideos(context.Background())
	require.NoError(t, err)
	assert.Len(t, dbVideos, len(videos))
}

func TestClearAllTables_Integration(t *testing.T) {
	resetDB(t)

	service := NewAutoUploadDataService(querier)

	// Добавляем данные
	regions := []domain.Region{
		{
			Id: 1,
			Name: domain.Translation{
				ArmName: "Շիրակ",
				EngName: "Shirak",
				RuName:  "Ширак",
			},
		},
	}
	err := service.CreateRegions(context.Background(), regions)
	require.NoError(t, err)

	// Очищаем
	err = service.ClearAllTables(context.Background())
	require.NoError(t, err)

	// Проверяем, что таблицы действительно пусты
	dbRegions, err := querier.GetRegions(context.Background())
	require.NoError(t, err)
	assert.Empty(t, dbRegions)
}
