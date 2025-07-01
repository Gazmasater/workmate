package main

import (
	"log"
	"net/http"
	"workmate/internal/delivery/_http"
	"workmate/repository/memory"

	"workmate/usecase"
)

func main() {
	repo := memory.NewInMemoryRepo()
	uc := usecase.NewTaskUseCase(repo)
	handler := _http.NewHandler(uc)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler.Routes()))
}
