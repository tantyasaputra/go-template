package app

import (
	"context"
	"go-template/internal/config"
	"go-template/internal/database"
	"go-template/internal/log"
	"go-template/internal/service/example"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	internalHttp "go-template/internal/http"

	examplePG "go-template/internal/repository/example/postgres"

	"github.com/gin-contrib/cors"
)

// Run is entry point
func Run(cfg *config.Config) int {
	// Initiate Data Handler (GORM)
	var dataHandler database.DataHandler = database.NewDataHandler(cfg.DB)
	// Initiate Repository
	exampleRepo := examplePG.NewExampleRepository(dataHandler)
	// Initiate Service
	exampleService := example.NewExampleService(exampleRepo)

	engine := internalHttp.NewServer(
		dataHandler,
		&internalHttp.Services{
			ExampleService: exampleService,
		},
	).Build()

	// Configure CORS using gin middleware
	corsCfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300 * time.Second,
	}
	if strings.EqualFold(cfg.DEVMODE, "true") {
		corsCfg.AllowAllOrigins = true
	} else {
		corsCfg.AllowOriginFunc = func(origin string) bool {
			return strings.HasSuffix(origin, ".internal.xfers.com") || strings.HasSuffix(origin, ".internal.fazz.com")
		}
	}
	engine.Use(cors.New(corsCfg))

	s := &http.Server{
		Addr: "0.0.0.0:" + cfg.PORT,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 600,
		ReadTimeout:  time.Second * 600,
		IdleTimeout:  time.Second * 600,
		Handler:      engine,
	}
	go func() {
		log.Infow("Http Listening Initiated", "event", "http init", "port", cfg.PORT)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("Http Listening Error", "event", "http init", "error", err, "port", cfg.PORT)
		}
	}()

	// graceful reload
	reload := make(chan os.Signal, 1)
	signal.Notify(reload, syscall.SIGHUP)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	for {
		select {
		case <-reload:
			log.Infow("Http Restart Initiated", "event", "http restart", "service", "IKN_B2B")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			if err := s.Shutdown(ctx); err != nil {
				log.Fatalw("Http Restart Error", "http restart", "service", "IKN_B2B", "error", err)
			}

			cancel()
			return 1

		case <-quit:
			log.Infow("Http Shutdown Initiated", "event", "http shutdown", "service", "IKN_B2B")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			if err := s.Shutdown(ctx); err != nil {
				log.Fatalw("Http Shutdown Error", "http shutdown", "service", "IKN_B2B", "error", err)
			}
			log.Infow("Http Shutdown Completed", "event", "http shutdown", "service", "IKN_B2B")

			cancel()
			return 0
		}
	}
}
