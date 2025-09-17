package http

import (
	"go-template/internal/log"
	"strings"

	"github.com/gin-gonic/gin"
)

const requestIDKey = "REQUEST_ID"

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			requestID, ok := c.Request.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}

			if strings.Contains(c.Request.URL.Path, "/health") {
				return
			}

			log.Infow("Http Traffic",
				"request id", requestID,
				"http method", c.Request.Method,
				"url fragment", c.Request.URL.Path,
				"client address", c.ClientIP(),
				"user agent", c.Request.UserAgent(),
				"event", "request received")
		}()

		c.Next()
	}
}
