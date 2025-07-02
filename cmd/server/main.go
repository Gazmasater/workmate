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

	"github.com/gaz358/myprog/workmate/internal/delivery/phttp"
	"github.com/gaz358/myprog/workmate/pkg/logger"
	"github.com/gaz358/myprog/workmate/repository/memory"
	"github.com/gaz358/myprog/workmate/usecase"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/gaz358/myprog/workmate/cmd/server/docs"

	"github.com/go-chi/chi/v5"
)

func main() {
	logger.SetLevel(logger.InfoLevel)
	logg := logger.Global().Named("main")

	repo := memory.NewInMemoryRepo()
	uc := usecase.NewTaskUseCase(repo)
	handler := phttp.NewHandler(uc)

	r := chi.NewRouter()

	r.Mount("/tasks", handler.Routes())

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logg.Fatalw("Server forced to shutdown", "error", err)
	}
	logg.Infow("Server exited gracefully")
}
