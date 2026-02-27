package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	api "github.com/Ari-Pari/backend/internal/api/generated"
	"github.com/Ari-Pari/backend/internal/clients/filestorage"
	db "github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/Ari-Pari/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Server struct {
	logger  *log.Logger
	db      db.Querier              // бд
	storage filestorage.FileStorage // minio
	// Добавьте ваши зависимости (БД, кэш, сервисы и т.д.)
}

func (s *Server) PostDancesSearch(w http.ResponseWriter, r *http.Request, params api.PostDancesSearchParams) {

	var req api.DanceSearchRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page := 1
	size := 20
	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}
	if params.Size != nil && *params.Size > 0 {
		size = *params.Size
	}

	lang := "en"
	if params.Lang != nil && *params.Lang != "" {
		lang = *params.Lang
	}

	searchText := req.SearchText

	genresIn := make([]string, 0, len(req.Genres))
	for _, g := range req.Genres {
		genresIn = append(genresIn, string(g))
	}

	regionIdsIn := make([]int64, 0, len(req.Regions))
	for _, id := range req.Regions {
		regionIdsIn = append(regionIdsIn, int64(id))
	}

	complexitiesIn := make([]int32, 0, len(req.Complexities))
	for _, c := range req.Complexities {
		complexitiesIn = append(complexitiesIn, int32(c))
	}

	pacesIn := make([]int32, 0, len(req.Paces))
	for _, p := range req.Paces {
		pacesIn = append(pacesIn, int32(p))
	}

	gendersIn := make([]string, 0, len(req.Genders))
	for _, g := range req.Genders {
		gendersIn = append(gendersIn, string(g))
	}

	handshakesIn := make([]string, 0, len(req.Handshakes))
	for _, h := range req.Handshakes {
		handshakesIn = append(handshakesIn, string(h))
	}

	orderByPopularity := false
	orderByAlphabet := false // Соответствует OrderByName в SQLC
	orderByCreatedAt := false

	switch req.SortedBy {
	case api.Popularity:
		orderByPopularity = true
	case api.Alphabet:
		orderByAlphabet = true
	case api.CreatedBy:
		orderByCreatedAt = true
	}

	reverseOrder := false // Соответствует DESC если true
	if strings.ToUpper(string(req.SortType)) == "DESC" {
		reverseOrder = true
	}

	dbParams := db.SearchDancesParams{
		Lang:              lang,
		SearchText:        searchText,
		GenresIn:          genresIn,
		RegionIdsIn:       regionIdsIn,
		ComplexitiesIn:    complexitiesIn,
		GendersIn:         gendersIn,
		PacesIn:           pacesIn,
		HandshakesIn:      handshakesIn,
		OrderByPopularity: orderByPopularity,
		OrderByCreatedAt:  orderByCreatedAt,
		OrderByName:       orderByAlphabet,
		ReverseOrder:      reverseOrder,
		Limit:             int32(size),
		Offset:            int32((page - 1) * size),
	}

	rows, err := s.db.SearchDances(r.Context(), dbParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]api.DanceShortResponse, 0, len(rows))
	for _, d := range rows {
		id := int(d.ID)

		genderEnum := api.DanceShortResponseGender(d.Gender)

		genres := make([]api.Genre, 0, len(d.Genres))
		for _, g := range d.Genres {
			genres = append(genres, api.Genre(g))
		}

		handshakes := make([]api.Handshake, 0, len(d.Handshakes))
		for _, h := range d.Handshakes {
			handshakes = append(handshakes, api.Handshake(h))
		}

		paces := make([]int, 0, len(d.Paces))
		for _, p := range d.Paces {
			paces = append(paces, int(p))
		}

		var regions []api.RegionResponse

		if len(d.RegionIds) > 0 {
			regions = make([]api.RegionResponse, 0, len(d.RegionIds))
			for i := range d.RegionIds {
				if i < len(d.RegionNames) {
					idVal := int(d.RegionIds[i])
					nameVal := d.RegionNames[i]
					regions = append(regions, api.RegionResponse{
						Id:   idVal,
						Name: nameVal,
					})
				}
			}
		}

		photoURL := ""
		if d.PhotoLink.Valid && d.PhotoLink.String != "" {
			url, err := s.storage.GetFileURL(ctx, d.PhotoLink.String, time.Hour)
			if err == nil {
				photoURL = url
			}
		}

		resp = append(resp, api.DanceShortResponse{
			Id:         &id,
			Name:       d.Name,
			Complexity: int(d.Complexity.Int32),
			Gender:     genderEnum,
			Genres:     genres,
			Handshakes: handshakes,
			Paces:      paces,
			PhotoLink:  photoURL,
			Regions:    regions,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetDancesId(w http.ResponseWriter, r *http.Request, id int, params api.GetDancesIdParams) {
	ctx := r.Context()
	danceID := int64(id)

	var argLang pgtype.Text
	if params.Lang != nil {
		argLang = pgtype.Text{String: *params.Lang, Valid: true}
	}

	dbDance, err := s.db.GetDanceByID(ctx, db.GetDanceByIDParams{ID: danceID, Lang: argLang})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			s.logger.Printf("db error (dance): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	dbRegions, err := s.db.GetRegionsByDanceID(ctx, db.GetRegionsByDanceIDParams{
		DanceID: danceID,
		Lang:    argLang,
	})
	if err != nil {
		s.logger.Printf("db error (regions): %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbVideos, err := s.db.GetVideosByDanceID(ctx, danceID)
	if err != nil {
		s.logger.Printf("db error (videos): %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbSongs, err := s.db.GetSongsByDanceID(ctx, danceID)
	if err != nil {
		s.logger.Printf("db error (songs): %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	idVal := int(dbDance.ID)
	res := api.DanceFullResponse{
		Id:         &idVal,
		Name:       dbDance.Name,
		Complexity: int(dbDance.Complexity.Int32),
		Gender:     api.DanceFullResponseGender(dbDance.Gender),
		PhotoLink:  "",
	}

	if dbDance.PhotoKey.Valid && dbDance.PhotoKey.String != "" {
		res.PhotoLink, _ = s.storage.GetFileURL(ctx, dbDance.PhotoKey.String, time.Hour)
	}

	res.Paces = make([]int, len(dbDance.Paces))
	for i, p := range dbDance.Paces {
		res.Paces[i] = int(p)
	}

	res.Regions = make([]api.RegionResponse, len(dbRegions))
	for i, reg := range dbRegions {
		res.Regions[i] = api.RegionResponse{Id: int(reg.ID), Name: reg.Name}
	}

	res.Songs = make([]api.SongResponse, len(dbSongs))
	for i, song := range dbSongs {
		songLink := ""
		if song.FileKey != "" {
			songLink, _ = s.storage.GetFileURL(ctx, song.FileKey, time.Hour)
		}
		res.Songs[i] = api.SongResponse{
			Id:        int(song.ID),
			Name:      song.Name,
			Link:      songLink,
			Ensembles: []api.EnsembleResponse{}, // Если ансамблей пока нет, отдаем пустой массив
		}
	}

	src, les, perf := []api.VideoResponse{}, []api.VideoResponse{}, []api.VideoResponse{}
	for _, v := range dbVideos {
		vid := api.VideoResponse{
			Id:   int(v.ID),
			Name: v.Name,
			Link: v.Link,
		}

		switch domain.VideoType(strings.ToUpper(v.Type)) {
		case domain.Source:
			src = append(src, vid)
		case domain.Lesson:
			les = append(les, vid)
		case domain.Video:
			perf = append(perf, vid)
		}
	}
	res.SourceVideos, res.LessonVideos, res.PerformanceVideos = &src, &les, &perf

	if len(dbDance.Genres) > 0 {
		res.Genres = api.Genre(dbDance.Genres[0])
	}
	res.Handshakes = make([]api.Handshake, len(dbDance.Handshakes))
	for i, h := range dbDance.Handshakes {
		res.Handshakes[i] = api.Handshake(h)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.logger.Printf("json encode error: %v", err)
		w.WriteHeader(500)
	}
}

func (s *Server) GetRegions(w http.ResponseWriter, r *http.Request, params api.GetRegionsParams) {
	ctx := r.Context()

	var argLang pgtype.Text
	if params.Lang != nil {
		argLang = pgtype.Text{String: *params.Lang, Valid: true}
	}

	dbRegions, err := s.db.ListRegions(ctx, argLang)
	if err != nil {
		s.logger.Printf("failed to list regions: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := make(api.RegionListResponse, len(dbRegions))
	for i, reg := range dbRegions {
		response[i] = api.RegionResponse{
			Id:   int(reg.ID),
			Name: reg.Name,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Printf("json encode error: %v", err)
	}
}

func NewServer(logger *log.Logger, db db.Querier, storage filestorage.FileStorage) *Server {
	return &Server{
		logger:  logger,
		db:      db,
		storage: storage,
	}
}
