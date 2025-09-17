package http

import (
	"github.com/alexliesenfeld/health"
	"github.com/gin-gonic/gin"
)

// list of all routes and their handler
func (s *Server) routes() {
	// health checking
	s.engine.GET("/health", gin.WrapH(health.NewHandler(s.healthChecker)))

	s.engine.GET("/sample", s.example())
}
