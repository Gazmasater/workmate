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

	"workmate/internal/delivery/phttp"
	"workmate/pkg/logger"
	"workmate/repository/memory"
	"workmate/usecase"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "workmate/cmd/server/docs"
	//_ "workmate/docs"

	"github.com/go-chi/chi/v5"
)

func main() {
	// 1) Логгер
	logger.SetLevel(logger.InfoLevel)
	logg := logger.Global().Named("main")

	// 2) Репозиторий, юзкейс, handler
	repo := memory.NewInMemoryRepo()
	uc := usecase.NewTaskUseCase(repo)
	handler := phttp.NewHandler(uc)

	// 3) Создаём корневой chi.Router и монтируем в него:
	r := chi.NewRouter()

	// 3.1) Ваши маршруты API
	r.Mount("/tasks", handler.Routes())

	// 3.2) Swagger UI по пути /swagger/*
	//     httpSwagger.WrapHandler сам отдаёт статические файлы и индекс
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// 4) Конфигурируем и запускаем HTTP-сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// graceful shutdown
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
