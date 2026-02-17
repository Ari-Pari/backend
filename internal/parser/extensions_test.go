package parser

import (
	"testing"

	"github.com/Ari-Pari/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

// -------------------------------
// Тесты для toDomainTranslation
// -------------------------------

func TestToDomainTranslation(t *testing.T) {
	dto := NameDto{
		ArmName: "Անուն",
		EngName: "Name",
		RuName:  "Имя",
	}

	translation := toDomainTranslation(dto)

	assert.Equal(t, dto.ArmName, translation.ArmName)
	assert.Equal(t, dto.EngName, translation.EngName)
	assert.Equal(t, dto.RuName, translation.RuName)
}

// -------------------------------
// Тесты для toDomainGenre
// -------------------------------

func TestToDomainGenre_ValidGenres(t *testing.T) {
	tests := []struct {
		name     string
		dto      GenreDto
		expected domain.Genre
	}{
		{"War", War, domain.War},
		{"Road", Road, domain.Road},
		{"Cult", Cult, domain.Cult},
		{"Lyrical", Lyrical, domain.Lyrical},
		{"Reverse", Reverse, domain.Reverse},
		{"Ritual", Ritual, domain.Ritual},
		{"Community", Community, domain.Community},
		{"Hunting", Hunting, domain.Hunting},
		{"Pilgrimage", Pilgrimage, domain.Pilgrimage},
		{"Memorable", Memorable, domain.Memorable},
		{"Memorial", Memorial, domain.Memorial},
		{"Funeral", Funeral, domain.Funeral},
		{"Festive", Festive, domain.Festive},
		{"Wedding", Wedding, domain.Wedding},
		{"Matchmakers", Matchmakers, domain.Matchmakers},
		{"Labor", Labor, domain.Labor},
		{"Amulet", Amulet, domain.Amulet},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, toDomainGenre(tt.dto))
		})
	}
}

func TestToDomainGenre_UnknownGenre(t *testing.T) {
	assert.Equal(t, domain.Genre(""), toDomainGenre("UnknownGenre"))
}

// -------------------------------
// Тесты для toDomainGender
// -------------------------------

func TestToDomainGender(t *testing.T) {
	assert.Equal(t, domain.Male, toDomainGender(Male))
	assert.Equal(t, domain.Female, toDomainGender(Female))
	assert.Equal(t, domain.Multi, toDomainGender(Multi))
	assert.Equal(t, domain.Gender(""), toDomainGender("Unknown"))
}

// -------------------------------
// Тесты для toDomainHoldingType
// -------------------------------

func TestToDomainHoldingType(t *testing.T) {
	assert.Equal(t, domain.Free, toDomainHoldingType(Free))
	assert.Equal(t, domain.LittleFinger, toDomainHoldingType(LittleFinger))
	assert.Equal(t, domain.Palm, toDomainHoldingType(Palm))
	assert.Equal(t, domain.Crossed, toDomainHoldingType(Crossed))
	assert.Equal(t, domain.Back, toDomainHoldingType(Back))
	assert.Equal(t, domain.Belt, toDomainHoldingType(Belt))
	assert.Equal(t, domain.Shoulder, toDomainHoldingType(Shoulder))
	assert.Equal(t, domain.Dagger, toDomainHoldingType(Dagger))
	assert.Equal(t, domain.Whip, toDomainHoldingType(Whip))
	assert.Equal(t, domain.HoldingType(""), toDomainHoldingType("Unknown"))
}

// -------------------------------
// Тесты для ToDomainVideoType
// -------------------------------

func TestToDomainVideoType(t *testing.T) {
	assert.Equal(t, domain.Lesson, ToDomainVideoType(VideoTypeLesson))
	assert.Equal(t, domain.Video, ToDomainVideoType(VideoTypeVideo))
	assert.Equal(t, domain.Source, ToDomainVideoType(VideoTypeSource))
	assert.Equal(t, domain.VideoType(""), ToDomainVideoType("Unknown"))
}

// -------------------------------
// Тесты для toDomainDance
// -------------------------------

func TestToDomainDance(t *testing.T) {
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

	domainDance := toDomainDance(dto)

	assert.Equal(t, dto.Id, domainDance.Id)
	assert.Equal(t, dto.NameKey, domainDance.NameKey)
	assert.Equal(t, toDomainTranslation(dto.Name), domainDance.Name)
	assert.Equal(t, dto.Difficult, domainDance.Complexity)
	assert.Equal(t, domain.Male, domainDance.Gender)
	assert.Equal(t, dto.Temps, domainDance.Paces)
	assert.Equal(t, []domain.Genre{domain.War, domain.Ritual}, domainDance.Genres)
	assert.Equal(t, []domain.HoldingType{domain.Free, domain.Belt}, domainDance.HoldingTypes)
	assert.Equal(t, dto.StateIds, domainDance.RegionIds)
	assert.Nil(t, domainDance.DeletedAt)
}

func TestToDomainDance_WithDeletedAt(t *testing.T) {
	difficult := int32(1)

	dto := DanceDto{
		Id:           456,
		NameKey:      "dance.extra",
		Name:         NameDto{ArmName: "Այլ", EngName: "Extra", RuName: "Другой"},
		Difficult:    &difficult,
		Gender:       Multi,
		Temps:        []int32{7, 8, 9},
		Genres:       []GenreDto{Community},
		HoldingTypes: []HoldingTypeDto{Dagger},
		StateIds:     []int64{4, 5},
		Type:         Extra,
	}

	domainDance := toDomainDance(dto)

	assert.NotNil(t, domainDance.DeletedAt)
	assert.NotZero(t, *domainDance.DeletedAt)
}

// -------------------------------
// Тесты для toDomainSong
// -------------------------------

func TestToDomainSong(t *testing.T) {
	dto := MusicDto{
		Id:       789,
		Name:     NameDto{ArmName: "Երգ", EngName: "Song", RuName: "Песня"},
		NameKey:  "song.key",
		DanceIds: []int64{101, 102},
	}

	domainSong := toDomainSong(dto)

	assert.Equal(t, dto.Id, domainSong.Id)
	assert.Equal(t, toDomainTranslation(dto.Name), domainSong.Name)
	assert.Equal(t, dto.NameKey, domainSong.NameKey)
	assert.Equal(t, dto.DanceIds, domainSong.DanceIds)
}

// -------------------------------
// Тесты для toDomainRegion
// -------------------------------

func TestToDomainRegion(t *testing.T) {
	dto := StateDto{
		Id:   321,
		Name: NameDto{ArmName: "Մարզ", EngName: "Region", RuName: "Регион"},
	}

	domainRegion := toDomainRegion(dto)

	assert.Equal(t, dto.Id, domainRegion.Id)
	assert.Equal(t, toDomainTranslation(dto.Name), domainRegion.Name)
}

// -------------------------------
// Тесты для toDomainVideo
// -------------------------------

func TestToDomainVideo(t *testing.T) {
	dto := VideoDto{
		Name:     NameDto{ArmName: "Տեսանյութ", EngName: "Video", RuName: "Видео"},
		Url:      "https://example.com/video",
		Type:     VideoTypeVideo,
		DanceIds: []int64{201, 202},
	}

	domainVideo := toDomainVideo(dto)

	assert.Nil(t, domainVideo.Id)
	assert.Equal(t, toDomainTranslation(dto.Name), domainVideo.Name)
	assert.Equal(t, dto.Name.ArmName, domainVideo.NameKey)
	assert.Equal(t, dto.Url, domainVideo.Link)
	assert.Equal(t, dto.DanceIds, domainVideo.DanceIds)
	assert.Equal(t, domain.Video, domainVideo.Type)
}

// -------------------------------
// Тесты для списков: ToDomainRegions, ToDomainDances и т.д.
// -------------------------------

func TestToDomainRegions(t *testing.T) {
	dtos := []StateDto{
		{Id: 1, Name: NameDto{ArmName: "Մարզ 1"}},
		{Id: 2, Name: NameDto{ArmName: "Մարզ 2"}},
	}

	regions := ToDomainRegions(dtos)
	assert.Len(t, regions, 2)
	assert.Equal(t, dtos[0].Id, regions[0].Id)
	assert.Equal(t, dtos[0].Name.ArmName, regions[0].Name.ArmName)
}

func TestToDomainDances(t *testing.T) {
	dtos := []DanceDto{
		{Id: 1, Name: NameDto{ArmName: "Պար 1"}},
		{Id: 2, Name: NameDto{ArmName: "Պար 2"}},
	}

	dances := ToDomainDances(dtos)
	assert.Len(t, dances, 2)
	assert.Equal(t, dtos[0].Id, dances[0].Id)
	assert.Equal(t, dtos[0].Name.ArmName, dances[0].Name.ArmName)
}

func TestToDomainSongs(t *testing.T) {
	dtos := []MusicDto{
		{Id: 1, Name: NameDto{ArmName: "Երգ 1"}},
		{Id: 2, Name: NameDto{ArmName: "Երգ 2"}},
	}

	songs := ToDomainSongs(dtos)
	assert.Len(t, songs, 2)
	assert.Equal(t, dtos[0].Id, songs[0].Id)
	assert.Equal(t, dtos[0].Name.ArmName, songs[0].Name.ArmName)
}

func TestToDomainVideos(t *testing.T) {
	dtos := []VideoDto{
		{Name: NameDto{ArmName: "Տեսանյութ 1"}},
		{Name: NameDto{ArmName: "Տեսանյութ 2"}},
	}

	videos := ToDomainVideos(dtos)
	assert.Len(t, videos, 2)
	assert.Equal(t, dtos[0].Name.ArmName, videos[0].Name.ArmName)
}
