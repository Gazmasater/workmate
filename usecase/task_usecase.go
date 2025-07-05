package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/gaz358/myprog/workmate/domain"
	"github.com/google/uuid"
)

type TaskUseCase struct {
	repo      domain.TaskRepository
	duration  time.Duration
	cancelMap map[string]context.CancelFunc
	mu        sync.Mutex
}

func NewTaskUseCase(repo domain.TaskRepository, duration time.Duration) *TaskUseCase {
	return &TaskUseCase{
		repo:      repo,
		duration:  duration,
		cancelMap: make(map[string]context.CancelFunc),
	}
}

func (uc *TaskUseCase) CreateTask() (*domain.Task, error) {
	ctx := context.Background() // Можно объявить в начале теста, если его ещё нет

	task := &domain.Task{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		Status:    domain.StatusPending,
	}
	if err := uc.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	uc.mu.Lock()
	uc.cancelMap[task.ID] = cancel
	uc.mu.Unlock()

	// Делаем копию задачи и передаём по указателю (чтобы не было race)
	copy := *task
	go uc.run(ctx, &copy)
	return task, nil
}

func (uc *TaskUseCase) run(ctx context.Context, task *domain.Task) {
	task.Status = domain.StatusRunning
	task.StartedAt = time.Now()
	_ = uc.repo.Update(ctx, task)

	select {
	case <-ctx.Done():
		task.Status = domain.StatusCancelled
		task.Result = "Canceled"
		task.EndedAt = time.Now()
		task.Duration = task.EndedAt.Sub(task.StartedAt).String()
		_ = uc.repo.Update(ctx, task)
	case <-time.After(uc.duration):
		task.Status = domain.StatusCompleted
		task.EndedAt = time.Now()
		task.Duration = task.EndedAt.Sub(task.StartedAt).String()
		task.Result = "OK"
		_ = uc.repo.Update(ctx, task)
	}

	// Чистим cancelMap
	uc.mu.Lock()
	delete(uc.cancelMap, task.ID)
	uc.mu.Unlock()
}

func (uc *TaskUseCase) GetTask(id string) (*domain.Task, error) {
	ctx := context.Background() // Можно объявить в начале теста, если его ещё нет

	return uc.repo.Get(ctx, id)
}

func (uc *TaskUseCase) DeleteTask(id string) error {
	ctx := context.Background() // Можно объявить в начале теста, если его ещё нет

	uc.mu.Lock()
	if cancel, ok := uc.cancelMap[id]; ok {
		cancel() // отменим если есть
		delete(uc.cancelMap, id)
	}
	uc.mu.Unlock()
	return uc.repo.Delete(ctx, id)
}

func (uc *TaskUseCase) ListTasks() ([]*domain.Task, error) {
	ctx := context.Background() // Можно объявить в начале теста, если его ещё нет

	return uc.repo.List(ctx)
}

func (uc *TaskUseCase) CancelTask(id string) error {
	uc.mu.Lock()
	cancel, ok := uc.cancelMap[id]
	uc.mu.Unlock()
	if !ok {
		return domain.ErrNotFound
	}
	cancel()
	return nil
}
