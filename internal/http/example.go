package http

import (
	"go-template/internal/log"
	"net/http"
)

func (s *Server) example() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.service.ExampleService.ExampleAdd(r.Context(), []string{
			"User 1",
			"User 2",
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		get, err := s.service.ExampleService.ExampleGet(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// get file from request
		// call parser using file
		// convert to standard format
		// concurrency -> call service
		log.Infow("Controller Info", "event", "service result", "message", get)
		w.WriteHeader(http.StatusOK)
	}
}
