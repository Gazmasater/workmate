package phttp

import (
	"encoding/json"
	"net/http"

	"workmate/pkg/logger"
	"workmate/usecase"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	uc  *usecase.TaskUseCase
	log logger.TypeOfLogger
}

func NewHandler(uc *usecase.TaskUseCase) *Handler {
	l := logger.Global().Named("http")
	return &Handler{
		uc:  uc,
		log: l,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Delete("/{id}", h.delete)
	return r
}

// @Summary      Создать новую задачу
// @Description  Инициализирует задачу со статусом Pending и возвращает её с сгенерированным ID
// @Tags         tasks
// @Produce      json

// @Router       /tasks [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	h.log.Infow("create task request", "method", r.Method, "path", r.URL.Path)

	task, err := h.uc.CreateTask()
	if err != nil {
		h.log.Errorw("failed to create task", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Infow("task created", "id", task.ID)
	writeJSON(w, task)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.log.Infow("get task request", "method", r.Method, "path", r.URL.Path, "id", id)

	task, err := h.uc.GetTask(id)
	if err != nil {
		h.log.Warnw("task not found", "id", id)
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	h.log.Infow("task retrieved", "id", task.ID)
	writeJSON(w, task)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.log.Infow("delete task request", "method", r.Method, "path", r.URL.Path, "id", id)

	if err := h.uc.DeleteTask(id); err != nil {
		h.log.Errorw("failed to delete task", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Infow("task deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
