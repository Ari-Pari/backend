package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	api "github.com/Ari-Pari/backend/internal/api/generated"
	"github.com/Ari-Pari/backend/internal/clients/filestorage"
	"github.com/Ari-Pari/backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Server struct {
	logger  *log.Logger
	db      db.Querier              // бд
	storage filestorage.FileStorage // minio
	// Добавьте ваши зависимости (БД, кэш, сервисы и т.д.)
}

func (s Server) PostDancesSearch(w http.ResponseWriter, r *http.Request, params api.PostDancesSearchParams) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetDancesId(w http.ResponseWriter, r *http.Request, id int, params api.GetDancesIdParams) {
	ctx := r.Context()

	var argLang pgtype.Text
	if params.Lang != nil {
		argLang = pgtype.Text{String: *params.Lang, Valid: true}
	}

	dbData, err := s.db.GetDanceByID(ctx, db.GetDanceByIDParams{
		ID:   int64(id),
		Lang: argLang,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			s.logger.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Получение ссылки в s3
	photoURL := ""
	if dbData.PhotoKey.Valid && dbData.PhotoKey.String != "" {
		url, err := s.storage.GetFileURL(ctx, dbData.PhotoKey.String, time.Hour)
		if err == nil {
			photoURL = url
		}
	}

	// Начинаем собирать ответ
	// Если поле в Swagger помечено как required, это будет обычный тип
	// Если нет (omitempty) — это будет указатель
	idVal := int(dbData.ID)

	res := api.DanceFullResponse{
		Id:         &idVal,
		Name:       dbData.Name,
		Complexity: int(dbData.Complexity.Int32),
		PhotoLink:  photoURL,
		Gender:     api.DanceFullResponseGender(dbData.Gender),
	}

	res.Paces = make([]int, len(dbData.Paces))
	for i, p := range dbData.Paces {
		res.Paces[i] = int(p)
	}

	if regionsBytes, ok := dbData.RegionsJson.([]byte); ok && len(regionsBytes) > 0 {
		json.Unmarshal(regionsBytes, &res.Regions)
	}

	type GetVideosRow struct {
		ID            int64       `json:"id"`
		Link          string      `json:"link"`
		TranslationID pgtype.Int8 `json:"translation_id"`
		Name          string      `json:"name"`
		Type          string      `json:"type"`
	}

	if videosBytes, ok := dbData.VideosJson.([]byte); ok && len(videosBytes) > 0 {
		var rawVideos []dbVideo
		if err := json.Unmarshal(videosBytes, &rawVideos); err == nil {
			
			sourceV := []api.VideoResponse{}
			lessonV := []api.VideoResponse{}
			perfV   := []api.VideoResponse{}

			for _, v := range rawVideos {
				vid := api.VideoResponse{
					Id:   int(v.ID),
					Name: v.Name,
					Link: v.Link,
				}

				switch v.Type {
				case "source":
					sourceV = append(sourceV, vid)
				case "lesson":
					lessonV = append(lessonV, vid)
				case "performance":
					perfV = append(perfV, vid)
				}
			}
			res.SourceVideos = &sourceV
			res.LessonVideos = &lessonV
			res.PerformanceVideos = &perfV
		}
	}

	if len(dbData.Genres) > 0 {
		res.Genres = api.Genre(dbData.Genres[0])
	}
	res.Handshakes = make([]api.Handshake, len(dbData.Handshakes))
	for i, h := range dbData.Handshakes {
		res.Handshakes[i] = api.Handshake(h)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}




func (s Server) GetRegions(w http.ResponseWriter, r *http.Request, params api.GetRegionsParams) {
	ctx := r.Context()

	// lang может не быть
	argLang := pgtype.Text{Valid: false}
	if params.Lang != nil {
		argLang = pgtype.Text{String: *params.Lang, Valid: true}
	}

	regions, err := s.db.ListRegions(ctx, argLang)
	if err != nil {
		s.logger.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(regions)
}

func NewServer(logger *log.Logger, db db.Querier, storage filestorage.FileStorage) *Server {
	return &Server{
		logger:  logger,
		db:      db,
		storage: storage,
	}
}
