package parser

import (
	"context"
	"os"
	"testing"

	"github.com/Ari-Pari/backend/internal/domain"
	"github.com/Ari-Pari/backend/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeFileReader — минимальная реализация FileReader для тестов
type fakeFileReader struct{}

func (f fakeFileReader) Open(name string) (*os.File, error) {
	// Создаём временный файл с фиктивным содержимым
	tmpfile, err := os.CreateTemp("", "fakefile-")
	if err != nil {
		return nil, err
	}
	_, err = tmpfile.Write([]byte("fake-content"))
	if err != nil {
		return nil, err
	}
	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	return tmpfile, nil
}

func (f fakeFileReader) ReadFile(name string) ([]byte, error) {
	return []byte("fake-content"), nil
}

// -------------------------------
// Тесты для toDomainDance
// -------------------------------

func TestToDomainDance_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockFileStorage(ctrl)

	difficult := int32(3)

	dto := DanceDto{
		Id:           123,
		NameKey:      "dance.key",
		Name:         NameDto{ArmName: "Գին", EngName: "Gin", RuName: "Гин"},
		Difficult:    &difficult,
		Gender:       Male,
		Temps:        []int32{4, 5, 6},
		Genres:       []GenreDto{War, Ritual},
		HoldingTypes: []HoldingTypeDto{Free, Belt},
		StateIds:     []int64{1, 2, 3},
		Type:         Active,
	}

	// Ожидаем вызов UploadFile
	mockStorage.EXPECT().
		UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("mock-image-key", nil).
		Times(1)

	ctx := context.Background()
	reader := fakeFileReader{}
	imageFolder := "/images/"

	domainDance, err := toDomainDance(ctx, mockStorage, reader, dto, imageFolder)

	require.NoError(t, err)
	require.NotNil(t, domainDance)

	assert.Equal(t, dto.Id, domainDance.Id)
	assert.Equal(t, dto.NameKey, domainDance.NameKey)
	assert.Equal(t, "Գին", domainDance.Name.ArmName)
	assert.Equal(t, dto.Difficult, domainDance.Complexity)
	assert.Equal(t, domain.Male, domainDance.Gender)
	assert.Equal(t, dto.Temps, domainDance.Paces)
	assert.Equal(t, []domain.Genre{domain.War, domain.Ritual}, domainDance.Genres)
	assert.Equal(t, []domain.HoldingType{domain.Free, domain.Belt}, domainDance.HoldingTypes)
	assert.Equal(t, dto.StateIds, domainDance.RegionIds)
	assert.Nil(t, domainDance.DeletedAt)

	require.NotNil(t, domainDance.FileKey)
	assert.Equal(t, "mock-image-key", *domainDance.FileKey)
}

func TestToDomainDance_UploadFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockFileStorage(ctrl)

	dto := DanceDto{
		Id:      456,
		Name:    NameDto{ArmName: "Այլ", EngName: "Extra", RuName: "Другой"},
		NameKey: "dance.extra",
		Type:    Extra,
		Genres:  []GenreDto{Community},
		Gender:  Multi,
		Temps:   []int32{7, 8, 9},
	}

	// Ожидаем ошибку при загрузке файла
	mockStorage.EXPECT().
		UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", assert.AnError).
		Times(1)

	ctx := context.Background()
	reader := fakeFileReader{}
	imageFolder := "/images/"

	domainDance, err := toDomainDance(ctx, mockStorage, reader, dto, imageFolder)

	require.Error(t, err)
	assert.Equal(t, domainDance, domain.DanceShort{})
}

// -------------------------------
// Тесты для toDomainSong
// -------------------------------

func TestToDomainSong_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockFileStorage(ctrl)

	dto := MusicDto{
		Id:       789,
		Name:     NameDto{ArmName: "Երգ", EngName: "Song", RuName: "Песня"},
		NameKey:  "song.key",
		DanceIds: []int64{101, 102},
	}

	// Ожидаем вызов UploadFile
	mockStorage.EXPECT().
		UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("mock-audio-key", nil).
		Times(1)

	ctx := context.Background()
	reader := fakeFileReader{}
	musicFolder := "/music/"

	domainSong, err := toDomainSong(ctx, mockStorage, reader, dto, musicFolder)

	require.NoError(t, err)
	require.NotNil(t, domainSong)

	assert.Equal(t, dto.Id, domainSong.Id)
	assert.Equal(t, dto.NameKey, domainSong.NameKey)
	assert.Equal(t, "Երգ", domainSong.Name.ArmName)
	assert.Equal(t, dto.DanceIds, domainSong.DanceIds)

	require.NotNil(t, domainSong.FileKey)
	assert.Equal(t, "mock-audio-key", *domainSong.FileKey)
}

func TestToDomainSong_UploadFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockFileStorage(ctrl)

	dto := MusicDto{
		Id:       789,
		Name:     NameDto{ArmName: "Երգ", EngName: "Song", RuName: "Песня"},
		NameKey:  "song.key",
		DanceIds: []int64{101, 102},
	}

	// Ожидаем ошибку при загрузке файла
	mockStorage.EXPECT().
		UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", assert.AnError).
		Times(1)

	ctx := context.Background()
	reader := fakeFileReader{}
	musicFolder := "/music/"

	domainSong, err := toDomainSong(ctx, mockStorage, reader, dto, musicFolder)

	require.Error(t, err)
	assert.Equal(t, domain.SongShort{}, domainSong)
}

// -------------------------------
// Тесты для списков
// -------------------------------

func TestToDomainDances_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockFileStorage(ctrl)

	dtos := []DanceDto{
		{
			Id:      1,
			Name:    NameDto{ArmName: "Պար 1"},
			NameKey: "dance1.key",
			Type:    Active,
		},
		{
			Id:      2,
			Name:    NameDto{ArmName: "Պար 2"},
			NameKey: "dance2.key",
			Type:    Extra,
		},
	}

	// Ожидаем два вызова UploadFile
	mockStorage.EXPECT().
		UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("mock-key-1", nil).
		Times(1)

	mockStorage.EXPECT().
		UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("mock-key-2", nil).
		Times(1)

	ctx := context.Background()
	reader := fakeFileReader{}
	imageFolder := "/images/"

	dances, err := ToDomainDances(ctx, mockStorage, reader, dtos, imageFolder)

	require.NoError(t, err)
	require.Len(t, dances, 2)

	assert.Equal(t, int64(1), dances[0].Id)
	assert.Equal(t, "Պար 1", dances[0].Name.ArmName)
	assert.Equal(t, "mock-key-1", *dances[0].FileKey)

	assert.Equal(t, int64(2), dances[1].Id)
	assert.Equal(t, "Պար 2", dances[1].Name.ArmName)
	assert.Equal(t, "mock-key-2", *dances[1].FileKey)
}

func TestToDomainSongs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockFileStorage(ctrl)

	dtos := []MusicDto{
		{
			Id:       1,
			Name:     NameDto{ArmName: "Երգ 1"},
			NameKey:  "song1.key",
			DanceIds: []int64{1, 2},
		},
		{
			Id:       2,
			Name:     NameDto{ArmName: "Երգ 2"},
			NameKey:  "song2.key",
			DanceIds: []int64{3, 4},
		},
	}

	// Ожидаем два вызова UploadFile
	mockStorage.EXPECT().
		UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("mock-song-key-1", nil).
		Times(1)

	mockStorage.EXPECT().
		UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("mock-song-key-2", nil).
		Times(1)

	ctx := context.Background()
	reader := fakeFileReader{}
	musicFolder := "/music/"

	songs := ToDomainSongs(ctx, mockStorage, reader, dtos, musicFolder)

	require.Len(t, songs, 2)

	assert.Equal(t, int64(1), songs[0].Id)
	assert.Equal(t, "Երգ 1", songs[0].Name.ArmName)
	assert.Equal(t, "mock-song-key-1", *songs[0].FileKey)

	assert.Equal(t, int64(2), songs[1].Id)
	assert.Equal(t, "Երգ 2", songs[1].Name.ArmName)
	assert.Equal(t, "mock-song-key-2", *songs[1].FileKey)
}
