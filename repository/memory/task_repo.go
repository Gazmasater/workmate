package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/gaz358/myprog/workmate/domain"
	"github.com/gaz358/myprog/workmate/pkg/logger"
)

type InMemoryRepo struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{tasks: make(map[string]*domain.Task)}
}

func (r *InMemoryRepo) Create(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; exists {
		logger.WarnKV(ctx, "[memory] Create: task already exists", "task_id", task.ID)
		return errors.New("task already exists")
	}
	tCopy := *task
	r.tasks[task.ID] = &tCopy
	logger.InfoKV(ctx, "[memory] Create: task created", "task_id", task.ID)
	return nil
}

func (r *InMemoryRepo) Update(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; !exists {
		logger.WarnKV(ctx, "[memory] Update: task not found", "task_id", task.ID)
		return domain.ErrNotFound
	}
	tCopy := *task
	r.tasks[task.ID] = &tCopy
	logger.InfoKV(ctx, "[memory] Update: task updated", "task_id", task.ID)
	return nil
}

func (r *InMemoryRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		logger.WarnKV(ctx, "[memory] Delete: task not found", "task_id", id)
		return domain.ErrNotFound
	}
	delete(r.tasks, id)
	logger.InfoKV(ctx, "[memory] Delete: task deleted", "task_id", id)
	return nil
}

func (r *InMemoryRepo) Get(ctx context.Context, id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		logger.WarnKV(ctx, "[memory] Get: task not found", "task_id", id)
		return nil, domain.ErrNotFound
	}
	tCopy := *t
	logger.InfoKV(ctx, "[memory] Get: task found", "task_id", id)
	return &tCopy, nil
}

func (r *InMemoryRepo) List(ctx context.Context) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tCopy := *task
		result = append(result, &tCopy)
	}
	logger.InfoKV(ctx, "[memory] List: all tasks listed", "count", len(result))
	return result, nil
}
