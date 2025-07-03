// @title           Tasks API
// @version         1.0
// @description     Сервис управления задачами
// @host            localhost:8080
// @BasePath        /
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gaz358/myprog/workmate/config"
	"github.com/gaz358/myprog/workmate/internal/delivery/phttp"
	"github.com/gaz358/myprog/workmate/pkg/logger"
	"github.com/gaz358/myprog/workmate/repository/memory"
	"github.com/gaz358/myprog/workmate/usecase"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/gaz358/myprog/workmate/cmd/server/docs"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	logger.SetLevel(parseLogLevel(cfg.LogLevel))
	logg := logger.Global().Named("main")

	repo := memory.NewInMemoryRepo()
	uc := usecase.NewTaskUseCase(repo, cfg.TaskDuration)
	handler := phttp.NewHandler(uc)

	r := chi.NewRouter()
	r.Mount("/tasks", handler.Routes())
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		logg.Infow("Starting HTTP server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatalw("ListenAndServe failed", "error", err)
		}
	}()

	<-quit
	logg.Infow("Shutting down server…")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logg.Fatalw("Server forced to shutdown", "error", err)
	}
	logg.Infow("Server exited gracefully")
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
