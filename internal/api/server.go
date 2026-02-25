package api

import (
	"encoding/json"
	api "github.com/Ari-Pari/backend/internal/api/generated"
	"github.com/Ari-Pari/backend/internal/clients/filestorage"
	db "github.com/Ari-Pari/backend/internal/db/sqlc"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	logger  *log.Logger
	queries *db.Queries
	storage filestorage.FileStorage // minio
	// Добавьте ваши зависимости (БД, кэш, сервисы и т.д.)
}

func (s Server) PostDancesSearch(w http.ResponseWriter, r *http.Request, params api.PostDancesSearchParams) {
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

	// === ИСПРАВЛЕНИЕ ОШИБОК С GENRES ===
	genresIn := make([]string, 0, len(req.Genres))
	for _, g := range req.Genres {
		genresIn = append(genresIn, string(g))
	}
	// === КОНЕЦ ИСПРАВЛЕНИЯ ОШИБОК С GENRES ===

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

	// === ИСПРАВЛЕНИЕ ОШИБОК С СОРТИРОВКОЙ ===
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
		OrderByName:       orderByAlphabet, // Соответствует OrderByAlphabet
		ReverseOrder:      reverseOrder,
		Limit:             int32(size),
		Offset:            int32((page - 1) * size),
	}

	rows, err := s.queries.SearchDances(r.Context(), dbParams)
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

		// === ИСПРАВЛЕНИЕ ОШИБОК С REGIONS (обновлено) ===
		var regions []api.RegionResponse

		// Приведение типов для RegionIds
		regionIDs, ok := d.RegionIds.([]int64)
		if !ok {
			s.logger.Printf("Warning: d.RegionIds is not []int64, got %T. Initializing as empty slice.", d.RegionIds)
			regionIDs = []int64{} // Если приведение не удалось, инициализируем пустым слайсом
		}

		// Приведение типов для RegionNames
		regionNames, ok := d.RegionNames.([]string)
		if !ok {
			s.logger.Printf("Warning: d.RegionNames is not []string, got %T. Initializing as empty slice.", d.RegionNames)
			regionNames = []string{} // Если приведение не удалось, инициализируем пустым слайсом
		}

		if len(regionIDs) > 0 {
			regions = make([]api.RegionResponse, 0, len(regionIDs))
			for i := range regionIDs {
				if i < len(regionNames) { // Убедимся, что имя существует для этого ID
					idVal := int(regionIDs[i])
					nameVal := regionNames[i]
					regions = append(regions, api.RegionResponse{
						Id:   idVal,
						Name: nameVal,
					})
				}
			}
		}
		// === КОНЕЦ ИСПРАВЛЕНИЯ ОШИБОК С REGIONS ===

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

func (s Server) GetDancesId(w http.ResponseWriter, r *http.Request, id int, params api.GetDancesIdParams) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetRegions(w http.ResponseWriter, r *http.Request, params api.GetRegionsParams) {
	//TODO implement me
	panic("implement me")
}

func NewServer(logger *log.Logger, q *db.Queries, storage filestorage.FileStorage) *Server {
	return &Server{
		logger:  logger,
		queries: q,
		storage: storage,
	}
}
