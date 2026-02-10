package api

import (
	"log"
	"net/http"

	api "github.com/Ari-Pari/backend/internal/api/generated"
)

type Server struct {
	logger *log.Logger
	// Добавьте ваши зависимости (БД, кэш, сервисы и т.д.)
}

func (s Server) PostDancesSearch(w http.ResponseWriter, r *http.Request, params api.PostDancesSearchParams) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetDancesId(w http.ResponseWriter, r *http.Request, id int, params api.GetDancesIdParams) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetGenres(w http.ResponseWriter, r *http.Request, params api.GetGenresParams) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetRegions(w http.ResponseWriter, r *http.Request, params api.GetRegionsParams) {
	//TODO implement me
	panic("implement me")
}

func NewServer(logger *log.Logger) *Server {
	return &Server{
		logger: logger,
	}
}
