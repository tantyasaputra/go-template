package http

import (
	"time"

	"go-template/internal/database"
	"go-template/internal/service/example"

	"github.com/alexliesenfeld/health"
	"github.com/gorilla/mux"
)

// Server is working as http handler
type Server struct {
	router        *mux.Router
	dataHandler   database.DataHandler
	service       *Services
	healthChecker health.Checker
}

// Services list all available service here for http layer to use
type Services struct {
	ExampleService example.Service
}

// NewServer initiate http handler
func NewServer(
	dh database.DataHandler,
	serv *Services,
) *Server {
	s := &Server{
		router:      mux.NewRouter(),
		dataHandler: dh,
		service:     serv,
	}

	s.healthChecker = s.buildHealth()

	s.router.Use(LoggingMiddleware())

	s.routes()
	return s
}

// Build return mux router
func (s *Server) Build() *mux.Router {
	return s.router
}

func (s *Server) buildHealth() health.Checker {
	checker := health.NewChecker(

		// Set the time-to-live for our cache to 1 second (default).
		health.WithCacheDuration(3*time.Second),

		// Configure a global timeout that will be applied to all checks.
		health.WithTimeout(10*time.Second),

		// A check configuration to see if our database connection is up.
		// The check function will be executed for each HTTP request.
		health.WithCheck(health.Check{
			Name:    "database",      // A unique check name.
			Timeout: 2 * time.Second, // A check specific timeout.
			Check:   s.dataHandler.Ping,
		}),
	)

	return checker
}
