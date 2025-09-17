package http

import (
	"net/http"

	"github.com/alexliesenfeld/health"
)

// list of all routes and their handler
func (s *Server) routes() {
	// health checking
	s.router.HandleFunc("/health", health.NewHandler(s.healthChecker)).Methods(http.MethodGet)

	s.router.HandleFunc("/sample", s.example()).Methods("GET")
}
