package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"workmate/internal/delivery/_http"
	"workmate/pkg/logger"
	"workmate/repository/memory"
	"workmate/usecase"

	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
)

func main() {
	// Настраиваем логгер
	logger.SetLevel(logger.InfoLevel)
	log := logger.Global().Named("main")

	// Создаём репозиторий, юзкейз и HTTP-хендлер
	repo := memory.NewInMemoryRepo()
	uc := usecase.NewTaskUseCase(repo)
	handler := _http.NewHandler(uc)

	// Конфигурируем HTTP-сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.Routes(),
	}

	// Канал для перехвата сигнала прерывания
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Запускаем сервер в горутине
	go func() {
		log.Infow("Starting HTTP server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("ListenAndServe failed", "error", err)
		}
	}()

	// Ждём SIGINT (CTRL+C)
	<-quit
	log.Infow("Shutting down server...")

	// Создаём контекст с таймаутом для завершения работы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalw("Server forced to shutdown", "error", err)
	}

	log.Infow("Server exited gracefully")
}
