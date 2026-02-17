package autoUploadDataService

import (
	"context"

	db "github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/Ari-Pari/backend/internal/domain"
)

type AutoUploadDataService interface {
	CreateVideos(ctx context.Context, videos []domain.VideoShort) error
	CreateRegions(ctx context.Context, regions []domain.Region) error
	CreateDances(ctx context.Context, dances []domain.DanceShort) error
	CreateSongs(ctx context.Context, songs []domain.SongShort) error
	ClearAllTables(ctx context.Context) error
}

type autoUploadDataService struct {
	querier db.Querier
}

func (a autoUploadDataService) CreateVideos(ctx context.Context, videos []domain.VideoShort) error {
	translations := make([]domain.Translation, len(videos))
	for i := range videos {
		translations[i] = videos[i].Name
	}

	translationToParams := TranslationToDao(translations)
	translationIds, err := a.querier.InsertTranslations(ctx, translationToParams)
	if err != nil {
		return err
	}

	videoParams := VideosToDao(videos, translationIds)
	videoIds, err := a.querier.InsertVideos(ctx, videoParams)
	if err != nil {
		return err
	}

	danceVideoToParams, err := DanceVideosToDao(videos, videoIds)

	if err != nil {
		return err
	}

	return a.querier.InsertDanceVideos(ctx, danceVideoToParams)
}

func (a autoUploadDataService) CreateSongs(ctx context.Context, songs []domain.SongShort) error {
	translations := make([]domain.Translation, len(songs))
	for i := range songs {
		translations[i] = songs[i].Name
	}
	translationToParams := TranslationToDao(translations)
	translationIds, err := a.querier.InsertTranslations(ctx, translationToParams)
	if err != nil {
		return err
	}
	songToParams := SongsToDao(songs, translationIds)
	if err = a.querier.InsertSongs(ctx, songToParams); err != nil {
		return err
	}

	danceSongToParams := SongDancesToDao(songs)
	return a.querier.InsertDanceSongs(ctx, danceSongToParams)
}

func (a autoUploadDataService) ClearAllTables(ctx context.Context) error {
	//TODO implement me
	return a.querier.TruncateAllTables(ctx)
}

func (a autoUploadDataService) CreateRegions(ctx context.Context, regions []domain.Region) error {
	translations := make([]domain.Translation, len(regions))
	for i := range regions {
		translations[i] = regions[i].Name
	}
	translationToParams := TranslationToDao(translations)
	translationIds, err := a.querier.InsertTranslations(ctx, translationToParams)
	if err != nil {
		return err
	}
	regionToParams := RegionToDao(regions, translationIds)
	return a.querier.InsertRegions(ctx, regionToParams)
}

func (a autoUploadDataService) CreateDances(ctx context.Context, dances []domain.DanceShort) error {
	translations := make([]domain.Translation, len(dances))
	for i := range dances {
		translations[i] = dances[i].Name
	}
	translationToParams := TranslationToDao(translations)
	translationIds, err := a.querier.InsertTranslations(ctx, translationToParams)
	if err != nil {
		return err
	}
	dancesToParams := DanceToDao(dances, translationIds)
	for _, dance := range dancesToParams {
		if err = a.querier.InsertDance(ctx, dance); err != nil {
			return err
		}
	}

	danceRegions := DanceRegionsToDao(dances)

	return a.querier.InsertDanceRegions(ctx, danceRegions)
}

func NewAutoUploadDataService(myQuerier db.Querier) AutoUploadDataService {
	return &autoUploadDataService{
		querier: myQuerier,
	}
}
