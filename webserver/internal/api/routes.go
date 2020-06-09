package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) registerRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/thresholds", s.getThresholdHandler()).Methods(http.MethodGet)
	r.Handle("/thresholds", s.postThresholdHandler()).Methods(http.MethodPost)

	return r
}
