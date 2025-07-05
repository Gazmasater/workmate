package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/gaz358/myprog/workmate/domain"
	"github.com/gaz358/myprog/workmate/pkg/logger" // Импорт логгера
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
	ctx := context.Background()
	task := &domain.Task{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		Status:    domain.StatusPending,
	}

	logger.InfoKV(ctx, "creating task", "task_id", task.ID)

	if err := uc.repo.Create(ctx, task); err != nil {
		logger.ErrorKV(ctx, "failed to create task", "task_id", task.ID, "err", err)
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	uc.mu.Lock()
	uc.cancelMap[task.ID] = cancel
	uc.mu.Unlock()

	copy := *task
	go uc.run(ctx, &copy)

	logger.InfoKV(ctx, "task created", "task_id", task.ID)
	return task, nil
}

func (uc *TaskUseCase) run(ctx context.Context, task *domain.Task) {
	logger.InfoKV(ctx, "running task", "task_id", task.ID)
	task.Status = domain.StatusRunning
	task.StartedAt = time.Now()
	if err := uc.repo.Update(ctx, task); err != nil {
		logger.ErrorKV(ctx, "failed to update task status to running", "task_id", task.ID, "err", err)
	}

	select {
	case <-ctx.Done():
		task.Status = domain.StatusCanceled
		task.Result = "Canceled"
		task.EndedAt = time.Now()
		task.Duration = task.EndedAt.Sub(task.StartedAt).String()
		if err := uc.repo.Update(ctx, task); err != nil {
			logger.ErrorKV(ctx, "failed to update canceled task", "task_id", task.ID, "err", err)
		}
		logger.InfoKV(ctx, "task canceled", "task_id", task.ID)
	case <-time.After(uc.duration):
		task.Status = domain.StatusCompleted
		task.EndedAt = time.Now()
		task.Duration = task.EndedAt.Sub(task.StartedAt).String()
		task.Result = "OK"
		if err := uc.repo.Update(ctx, task); err != nil {
			logger.ErrorKV(ctx, "failed to update completed task", "task_id", task.ID, "err", err)
		}
		logger.InfoKV(ctx, "task completed", "task_id", task.ID)
	}

	uc.mu.Lock()
	delete(uc.cancelMap, task.ID)
	uc.mu.Unlock()
	logger.DebugKV(ctx, "removed task from cancelMap", "task_id", task.ID)
}

func (uc *TaskUseCase) GetTask(id string) (*domain.Task, error) {
	ctx := context.Background()
	logger.DebugKV(ctx, "get task", "task_id", id)
	task, err := uc.repo.Get(ctx, id)
	if err != nil {
		logger.ErrorKV(ctx, "failed to get task", "task_id", id, "err", err)
	}
	return task, err
}

func (uc *TaskUseCase) DeleteTask(id string) error {
	ctx := context.Background()
	logger.InfoKV(ctx, "delete task", "task_id", id)
	uc.mu.Lock()
	if cancel, ok := uc.cancelMap[id]; ok {
		cancel()
		delete(uc.cancelMap, id)
		logger.InfoKV(ctx, "task canceled via delete", "task_id", id)
	}
	uc.mu.Unlock()
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		logger.ErrorKV(ctx, "failed to delete task", "task_id", id, "err", err)
	}
	return err
}

func (uc *TaskUseCase) ListTasks() ([]*domain.Task, error) {
	ctx := context.Background()
	logger.Debug(ctx, "list tasks")
	tasks, err := uc.repo.List(ctx)
	if err != nil {
		logger.ErrorKV(ctx, "failed to list tasks", "err", err)
	}
	return tasks, err
}

func (uc *TaskUseCase) CancelTask(id string) error {
	ctx := context.Background()
	uc.mu.Lock()
	cancel, ok := uc.cancelMap[id]
	uc.mu.Unlock()
	if !ok {
		logger.WarnKV(ctx, "cancel called but task not found in cancelMap", "task_id", id)
		return domain.ErrNotFound
	}
	logger.InfoKV(ctx, "cancel task", "task_id", id)
	cancel()
	return nil
}
