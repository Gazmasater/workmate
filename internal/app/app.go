package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gaz358/myprog/workmate/config"
	"github.com/gaz358/myprog/workmate/internal/delivery/health"
	"github.com/gaz358/myprog/workmate/internal/delivery/phttp"
	"github.com/gaz358/myprog/workmate/pkg/logger"
	"github.com/gaz358/myprog/workmate/repository/memory"
	"github.com/gaz358/myprog/workmate/usecase"

	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Точка входа для main.go
func Run() {
	cfg := config.Load()
	logg := initLogger(cfg)
	repo := memory.NewInMemoryRepo()
	uc := usecase.NewTaskUseCase(repo, cfg.TaskDuration)
	handler := phttp.NewHandler(uc)
	router := setupRouter(handler)
	server := newServer(cfg, router)
	runServer(server, logg, cfg.ShutdownTimeout)
}

func initLogger(cfg *config.Config) logger.TypeOfLogger {
	logger.SetLevel(parseLogLevel(cfg.LogLevel))
	return logger.Global().Named("main")
}

func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "debug":
		return logger.DebugLevel
	case "info":
		return logger.InfoLevel
	case "warn":
		return logger.WarnLevel
	case "error":
		return logger.ErrorLevel
	default:
		return logger.InfoLevel
	}
}

func setupRouter(handler *phttp.Handler) http.Handler {
	r := chi.NewRouter()
	r.Mount("/tasks", handler.Routes())
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/health", health.Handler)
	return r
}

func newServer(cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}
}

func runServer(srv *http.Server, logg logger.TypeOfLogger, shutdownTimeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		logg.Infow("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatalw("server error", "err", err)
		}
	}()

	<-quit
	logg.Infow("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logg.Fatalw("shutdown error", "err", err)
	}
	logg.Infow("server stopped")
}
