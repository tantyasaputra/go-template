package http

import (
	"go-template/internal/log"

	"github.com/gin-gonic/gin"
)

func (s *Server) example() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := s.service.ExampleService.ExampleAdd(c.Request.Context(), []string{
			"User 1",
			"User 2",
		})
		if err != nil {
			c.Status(500)
			return
		}
		get, err := s.service.ExampleService.ExampleGet(c.Request.Context())
		if err != nil {
			c.Status(500)
			return
		}
		log.Infow("Controller Info", "event", "service result", "message", get)
		c.Status(200)
	}
}
