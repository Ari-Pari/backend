package parser

import (
	"context"
	"strconv"
	"time"

	"github.com/Ari-Pari/backend/internal/clients/filestorage"
	"github.com/Ari-Pari/backend/internal/domain"
)

const AudioContentType string = "audio/mpeg"
const ImageContentType string = "image/jpeg"

func ToDomainRegions(states []StateDto) []domain.Region {
	regions := make([]domain.Region, len(states))
	for i, state := range states {
		regions[i] = toDomainRegion(state)
	}
	return regions
}

func ToDomainDances(ctx context.Context, storage filestorage.FileStorage, fileReader FileReader, dto []DanceDto, photosFolderName string) ([]domain.DanceShort, error) {
	dances := make([]domain.DanceShort, len(dto))
	for i, dance := range dto {
		var err error
		dances[i], err = toDomainDance(ctx, storage, fileReader, dance, photosFolderName)
		if err != nil {
			return []domain.DanceShort{}, err
		}
	}
	return dances, nil
}

func ToDomainSongs(ctx context.Context, storage filestorage.FileStorage, fileReader FileReader, dto []MusicDto, musicFolderName string) []domain.SongShort {
	songs := make([]domain.SongShort, len(dto))
	for i, song := range dto {
		var err error
		songs[i], err = toDomainSong(ctx, storage, fileReader, song, musicFolderName)
		if err != nil {
			return []domain.SongShort{}
		}
	}
	return songs
}

func ToDomainVideos(dto []VideoDto) []domain.VideoShort {
	videos := make([]domain.VideoShort, len(dto))
	for i, video := range dto {
		videos[i] = toDomainVideo(video)
	}
	return videos
}

func ToDomainArtists(dto []ArtistDto) []domain.ArtistShort {
	artists := make([]domain.ArtistShort, len(dto))
	for i, artist := range dto {
		artists[i] = toDomainArtist(artist)
	}
	return artists
}

func toDomainArtist(dto ArtistDto) domain.ArtistShort {
	now := time.Now()
	deletedAt := &now
	if dto.Type != Extra {
		deletedAt = nil
	}

	return domain.ArtistShort{
		Id:        dto.Id,
		Name:      toDomainTranslation(dto.Name),
		NameKey:   dto.Name.ArmName,
		Url:       dto.Url,
		DeletedAt: deletedAt,
	}
}

func toDomainVideo(dto VideoDto) domain.VideoShort {
	return domain.VideoShort{
		Id:       nil,
		Name:     toDomainTranslation(dto.Name),
		NameKey:  dto.Name.ArmName,
		Link:     dto.Url,
		DanceIds: dto.DanceIds,
		Type:     toDomainVideoType(dto.Type),
	}
}

func toDomainDance(ctx context.Context, storage filestorage.FileStorage, fileReader FileReader, dto DanceDto, imageFolderName string) (domain.DanceShort, error) {
	genres := make([]domain.Genre, len(dto.Genres))
	holdingTypes := make([]domain.HoldingType, len(dto.HoldingTypes))

	var deletedAt *time.Time = nil

	if dto.Type == Extra {
		deletedAt = new(time.Time)
		*deletedAt = time.Now()
	}

	for i, genre := range dto.Genres {
		genres[i] = toDomainGenre(genre)
	}

	for i, holdingType := range dto.HoldingTypes {
		holdingTypes[i] = toDomainHoldingType(holdingType)
	}

	key, err := parseFileForStorage(ctx, storage, fileReader, getImageFileName(imageFolderName, dto.Id), ImageContentType)

	if err != nil {
		return domain.DanceShort{}, err
	}

	return domain.DanceShort{
		Id:           dto.Id,
		Name:         toDomainTranslation(dto.Name),
		NameKey:      dto.NameKey,
		FileKey:      &key,
		Complexity:   dto.Difficult,
		Genres:       genres,
		Gender:       toDomainGender(dto.Gender),
		Paces:        dto.Temps,
		HoldingTypes: holdingTypes,
		RegionIds:    dto.StateIds,
		DeletedAt:    deletedAt,
	}, nil
}

func getImageFileName(folderName string, id int64) string {
	return folderName + strconv.FormatInt(id, 10) + ".jpeg"
}

func toDomainSong(ctx context.Context, storage filestorage.FileStorage, fileReader FileReader, dto MusicDto, musicFolderName string) (domain.SongShort, error) {
	fileKey, err := parseFileForStorage(ctx, storage, fileReader, getAudioFileName(musicFolderName, dto.Name.ArmName), AudioContentType)

	if err != nil {
		return domain.SongShort{}, err
	}

	return domain.SongShort{
		Id:        dto.Id,
		Name:      toDomainTranslation(dto.Name),
		NameKey:   dto.NameKey,
		FileKey:   &fileKey,
		DanceIds:  dto.DanceIds,
		ArtistIds: dto.Artists,
	}, nil
}

func getAudioFileName(folderName string, name string) string {
	return folderName + name + ".mp3"
}
func toDomainTranslation(dto NameDto) domain.Translation {
	return domain.Translation{
		ArmName: dto.ArmName,
		EngName: dto.EngName,
		RuName:  dto.RuName,
	}
}
func toDomainGenre(dto GenreDto) domain.Genre {
	switch dto {
	case War:
		return domain.War
	case Road:
		return domain.Road
	case Cult:
		return domain.Cult
	case Lyrical:
		return domain.Lyrical
	case Reverse:
		return domain.Reverse
	case Ritual:
		return domain.Ritual
	case Community:
		return domain.Community
	case Hunting:
		return domain.Hunting
	case Pilgrimage:
		return domain.Pilgrimage
	case Memorable:
		return domain.Memorable
	case Memorial:
		return domain.Memorial
	case Funeral:
		return domain.Funeral
	case Festive:
		return domain.Festive
	case Wedding:
		return domain.Wedding
	case Matchmakers:
		return domain.Matchmakers
	case Labor:
		return domain.Labor
	case Amulet:
		return domain.Amulet
	default:
		return ""
	}
}

func toDomainGender(dto GenderDto) domain.Gender {
	switch dto {
	case Male:
		return domain.Male
	case Female:
		return domain.Female
	case Multi:
		return domain.Multi
	default:
		return ""
	}
}

func toDomainHoldingType(dto HoldingTypeDto) domain.HoldingType {
	switch dto {
	case Free:
		return domain.Free
	case LittleFinger:
		return domain.LittleFinger
	case Palm:
		return domain.Palm
	case Crossed:
		return domain.Crossed
	case Back:
		return domain.Back
	case Belt:
		return domain.Belt
	case Shoulder:
		return domain.Shoulder
	case Dagger:
		return domain.Dagger
	case Whip:
		return domain.Whip
	default:
		return ""
	}
}

func toDomainRegion(dto StateDto) domain.Region {
	return domain.Region{
		Id:   dto.Id,
		Name: toDomainTranslation(dto.Name),
	}
}

func toDomainVideoType(dto VideoTypeDto) domain.VideoType {
	switch dto {
	case VideoTypeLesson:
		return domain.Lesson
	case VideoTypeVideo:
		return domain.Video
	case VideoTypeSource:
		return domain.Source
	default:
		return ""
	}
}
