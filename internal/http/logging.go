package http

import (
	"go-template/internal/log"
	"net/http"
	"strings"
)

const requestIDKey = "REQUEST_ID"

func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}

				if strings.Contains(r.URL.Path, "/health") {
					return
				}

				log.Infow("Http Traffic",
					"request id", requestID,
					"http method", r.Method,
					"url fragment", r.URL.Path,
					"client address", r.RemoteAddr,
					"user agent", r.UserAgent(),
					"event", "request received")
			}()

			next.ServeHTTP(w, r)
		})
	}
}
