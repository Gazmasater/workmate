package phttp

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gaz358/myprog/workmate/domain"
	"github.com/gaz358/myprog/workmate/pkg/logger"
	"github.com/gaz358/myprog/workmate/usecase"
	"github.com/go-chi/chi/v5"
)

type ErrorResponse struct {
	Message string `json:"message" example:"something went wrong"`
}

var _ = domain.Task{}

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
	r.Get("/all", h.list)

	r.Delete("/{id}", h.delete)
	r.Put("/{id}/cancel", h.cancel)
	r.Get("/health", h.Health) // health на корне API
	r.Get("/filter", h.filter)

	return r
}

// @Summary      Создать новую задачу
// @Description  Инициализирует задачу со статусом Pending и возвращает её с сгенерированным ID
// @Tags         tasks
// @Produce      json
// @Success      200  {object}  domain.Task         "Задача успешно создана"
// @Failure      500  {object}  ErrorResponse  "Внутренняя ошибка сервера"
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
	taskCopy := *task
	writeJSON(w, taskCopy)
}

// @Summary      Получить задачу по ID
// @Description  Возвращает задачу по её идентификатору
// @Tags         tasks
// @Produce      json
// @Param        id   path      string            true  "ID задачи"
// @Success      200  {object}  domain.Task        "Задача найдена"
// @Failure      404  {object}  phttp.ErrorResponse  "Задача не найдена"
// @Failure      500  {object}  phttp.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /tasks/{id} [get]
func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.log.Infow("get task request", "method", r.Method, "path", r.URL.Path, "id", id)

	task, err := h.uc.GetTask(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.log.Warnw("not found", "id", id)
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, ErrorResponse{Message: "not found"})
			return
		}

		h.log.Errorw("failed to get task", "id", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, ErrorResponse{Message: err.Error()})
		return
	}

	h.log.Infow("task retrieved", "id", task.ID)
	taskCopy := *task
	writeJSON(w, taskCopy)
}

// @Summary      Удалить задачу по ID
// @Description  Удаляет задачу из системы по её идентификатору
// @Tags         tasks
// @Param        id   path      string            true  "ID задачи"
// @Success      204  "No Content"
// @Failure      500  {object}  phttp.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /tasks/{id} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.log.Infow("delete task request", "method", r.Method, "path", r.URL.Path, "id", id)

	err := h.uc.DeleteTask(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.log.Warnw("not found", "id", id)
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, ErrorResponse{Message: "not found"})
			return
		}

		h.log.Errorw("failed to delete task", "id", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, ErrorResponse{Message: err.Error()})
		return
	}

	h.log.Infow("task deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// @Summary      Получить список всех задач
// @Tags         tasks
// @Produce      json
// @Success      200  {array}  domain.TaskListItem
// @Failure      500  {object}  ErrorResponse
// @Router       /tasks/all [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.uc.ListTasks()
	if err != nil {
		h.log.Errorw("failed to list tasks", "error", err)
		writeJSON(w, ErrorResponse{Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result []map[string]interface{}
	for _, t := range tasks {
		item := map[string]interface{}{
			"id":     t.ID,
			"status": t.Status,
		}
		if t.Status == domain.StatusCompleted {
			item["duration"] = t.EndedAt.Sub(t.StartedAt).String()
		}
		result = append(result, item)
	}

	writeJSON(w, result)
}

// @Summary      Отменить задачу
// @Description  Прерывает выполнение задачи, если она ещё не завершена
// @Tags         tasks
// @Param        id   path      string  true  "ID задачи"
// @Success      200  {object}  map[string]string  "Задача отменена"
// @Failure      404  {object}  ErrorResponse       "Задача не найдена"
// @Failure      500  {object}  ErrorResponse       "Внутренняя ошибка"
// @Router       /tasks/{id}/cancel [put]
func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.log.Infow("cancel task request", "method", r.Method, "path", r.URL.Path, "id", id)

	err := h.uc.CancelTask(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.log.Warnw("not found", "id", id)
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, ErrorResponse{Message: "not found"})
			return
		}
		h.log.Errorw("failed to cancel task", "id", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, ErrorResponse{Message: err.Error()})
		return
	}

	h.log.Infow("task canceled", "id", id)
	writeJSON(w, map[string]string{"status": "canceled"})
}

// @Summary      Healthcheck
// @Description  Проверка доступности сервиса
// @Tags         health
// @Produce      plain
// @Success      200 {string} string "ok"
// @Router       /health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Printf("error writing response: %v", err)
	}
}

// filter godoc
// @Summary      Фильтр и пагинация задач
// @Description  Фильтрует задачи по id, status, возвращает пагинацию
// @Tags         tasks
// @Produce      json
// @Param        id     query     string  false  "ID задачи (точное совпадение)"
// @Param        status query     string  false  "Статус задачи" Enums(pending, running, completed, failed, canceled)
// @Param        limit  query     int     false  "Максимум задач в ответе (default=10)"
// @Param        offset query     int     false  "Сдвиг (default=0)"
// @Success      200    {array}   domain.TaskListItem
// @Failure      500    {object}  ErrorResponse
// @Router       /tasks/filter [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")

	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	tasks, err := h.uc.ListTasks()
	if err != nil {
		h.log.Errorw("failed to list tasks", "error", err)
		writeJSON(w, ErrorResponse{Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Фильтрация
	filtered := make([]*domain.Task, 0)
	for _, t := range tasks {
		if id != "" && t.ID != id {
			continue
		}
		if status != "" && string(t.Status) != status {
			continue
		}
		filtered = append(filtered, t)
	}

	// Пагинация
	if offset > len(filtered) {
		offset = len(filtered)
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	paged := filtered[offset:end]

	writeJSON(w, paged)
}
