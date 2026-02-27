package autoUploadDataService

import (
	"time"

	db "github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/Ari-Pari/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func TranslationToDao(translations []domain.Translation) db.InsertTranslationsParams {

	engNames := make([]string, len(translations))
	ruNames := make([]string, len(translations))
	armNames := make([]string, len(translations))

	for i := range translations {
		engNames[i] = translations[i].EngName
		ruNames[i] = translations[i].RuName
		armNames[i] = translations[i].ArmName
	}

	return db.InsertTranslationsParams{
		EngNames: engNames,
		RuNames:  ruNames,
		ArmNames: armNames,
	}
}

func RegionToDao(regions []domain.Region, translationIds []int64) db.InsertRegionsParams {
	ids := make([]int64, len(regions))
	names := make([]string, len(regions))

	for i := range regions {
		ids[i] = regions[i].Id
		names[i] = regions[i].Name.ArmName
	}

	return db.InsertRegionsParams{
		Ids:            ids,
		TranslationIds: translationIds,
		Names:          names,
	}
}

func DanceToDao(dances []domain.DanceShort, translationIds []int64) []db.InsertDanceParams {
	params := make([]db.InsertDanceParams, len(dances))

	for i, dance := range dances {
		handshakes := make([]string, len(dance.HoldingTypes))
		for j, holdingType := range dance.HoldingTypes {
			handshakes[j] = string(holdingType)
		}

		genres := make([]string, len(dance.Genres))
		for j, genre := range dance.Genres {
			genres[j] = string(genre)
		}

		complexity := pgtype.Int4{Int32: 0, Valid: false}

		if dance.Complexity != nil {
			complexity = pgtype.Int4{Int32: *dance.Complexity, Valid: true}
		}

		deletedAt := pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: false,
		}

		if dance.DeletedAt != nil {
			deletedAt.Time = *dance.DeletedAt
			deletedAt.Valid = true
		}

		var photoKey pgtype.Text

		if dance.FileKey != nil {
			photoKey = pgtype.Text{
				String: *dance.FileKey,
				Valid:  true,
			}
		} else {
			photoKey = pgtype.Text{
				String: "",
				Valid:  false,
			}
		}

		params[i] = db.InsertDanceParams{
			ID: dance.Id,
			TranslationID: pgtype.Int8{
				Int64: translationIds[i],
				Valid: true,
			},
			Name:       dance.NameKey,
			PhotoKey:   photoKey,
			Paces:      dance.Paces,
			Gender:     string(dance.Gender),
			Complexity: complexity,
			Genres:     genres,
			DeletedAt:  deletedAt,
			Handshakes: handshakes,
			Popularity: 0,
		}
	}

	return params
}

func DanceRegionsToDao(dances []domain.DanceShort) db.InsertDanceRegionsParams {
	danceIds := make([]int64, 0, len(dances))
	regionIds := make([]int64, 0, len(dances))

	for _, dance := range dances {
		for _, region := range dance.RegionIds {
			danceIds = append(danceIds, dance.Id)
			regionIds = append(regionIds, region)
		}
	}

	return db.InsertDanceRegionsParams{
		DanceIds:  danceIds,
		RegionIds: regionIds,
	}
}

func SongsToDao(songs []domain.SongShort, translationIds []int64) db.InsertSongsParams {
	ids := make([]int64, len(songs))
	names := make([]string, len(songs))
	fileKeys := make([]string, len(songs))

	for i := range songs {
		ids[i] = songs[i].Id
		names[i] = songs[i].NameKey
		fileKeys[i] = *songs[i].FileKey
	}

	return db.InsertSongsParams{
		Ids:            ids,
		TranslationIds: translationIds,
		Names:          names,
		FileKeys:       fileKeys,
	}
}

func SongDancesToDao(songs []domain.SongShort) db.InsertDanceSongsParams {
	songIds := make([]int64, 0, len(songs))
	danceIds := make([]int64, 0, len(songs))

	for _, song := range songs {
		for _, dance := range song.DanceIds {
			songIds = append(songIds, song.Id)
			danceIds = append(danceIds, dance)
		}
	}

	return db.InsertDanceSongsParams{
		SongIds:  songIds,
		DanceIds: danceIds,
	}
}

func SongArtistsToDao(songs []domain.SongShort) db.InsertSongArtistsParams {
	songIds := make([]int64, 0, len(songs))
	artistIds := make([]int64, 0, len(songs))

	for _, song := range songs {
		for _, artistId := range song.ArtistIds {
			songIds = append(songIds, song.Id)
			artistIds = append(artistIds, artistId)
		}
	}

	return db.InsertSongArtistsParams{
		SongIds:   songIds,
		ArtistIds: artistIds,
	}
}

func VideosToDao(videos []domain.VideoShort, translationIds []int64) db.InsertVideosParams {
	links := make([]string, len(videos))
	names := make([]string, len(videos))
	types := make([]string, len(videos))

	for i := range videos {
		links[i] = videos[i].Link
		names[i] = videos[i].NameKey
		types[i] = string(videos[i].Type)
	}

	return db.InsertVideosParams{
		Links:          links,
		TranslationIds: translationIds,
		Names:          names,
		Types:          types,
	}
}

func DanceVideosToDao(videos []domain.VideoShort, videoIds []int64) db.InsertDanceVideosParams {
	if len(videos) != len(videoIds) {
		return db.InsertDanceVideosParams{}
	}
	danceIds := make([]int64, 0, len(videos))
	videoIdsToFill := make([]int64, 0, len(videos))

	for i, video := range videos {
		for _, dance := range video.DanceIds {
			danceIds = append(danceIds, dance)
			videoIdsToFill = append(videoIdsToFill, videoIds[i])
		}
	}

	return db.InsertDanceVideosParams{
		DanceIds: danceIds,
		VideoIds: videoIdsToFill,
	}
}

func ArtistsToDao(artists []domain.ArtistShort, translations []int64) db.InsertArtistsParams {
	ids := make([]int64, len(artists))
	names := make([]string, len(artists))
	links := make([]string, len(artists))
	deletedAts := make([]pgtype.Timestamptz, len(artists))

	for i, artist := range artists {
		ids[i] = artist.Id
		names[i] = artist.NameKey
		links[i] = artist.Url
		if artist.DeletedAt == nil {
			deletedAts[i] = pgtype.Timestamptz{Time: time.Now(), Valid: false}
		} else {
			deletedAts[i] = pgtype.Timestamptz{Time: *artist.DeletedAt, Valid: true}
		}
	}

	return db.InsertArtistsParams{
		Ids:            ids,
		TranslationIds: translations,
		Names:          names,
		Links:          links,
		DeletedAts:     deletedAts,
	}
}
