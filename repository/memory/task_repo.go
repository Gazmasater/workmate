package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/gaz358/myprog/workmate/domain"
	// logger не импортируем, если убрали все логи
)

type InMemoryRepo struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
	pool  sync.Pool
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		tasks: make(map[string]*domain.Task),
		pool: sync.Pool{
			New: func() interface{} {
				return new(domain.Task)
			},
		},
	}
}

func (r *InMemoryRepo) Create(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; exists {
		// return errors.New("task already exists") // Можно вернуть ошибку без лога
		return errors.New("task already exists")
	}
	tCopy := r.pool.Get().(*domain.Task)
	*tCopy = *task
	r.tasks[task.ID] = tCopy
	return nil
}

func (r *InMemoryRepo) Update(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; !exists {
		return domain.ErrNotFound
	}
	tCopy := r.pool.Get().(*domain.Task)
	*tCopy = *task
	// Старую версию с task.ID можешь вернуть в pool, если хочешь, вот так:
	// r.pool.Put(r.tasks[task.ID])
	r.tasks[task.ID] = tCopy
	return nil
}

func (r *InMemoryRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.tasks[id]
	if !ok {
		return domain.ErrNotFound
	}
	delete(r.tasks, id)
	r.pool.Put(t) // Возвращаем объект в pool!
	return nil
}

func (r *InMemoryRepo) Get(ctx context.Context, id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	// Возвращаем копию для безопасности
	tCopy := r.pool.Get().(*domain.Task)
	*tCopy = *t
	return tCopy, nil
}

func (r *InMemoryRepo) List(ctx context.Context) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tCopy := r.pool.Get().(*domain.Task)
		*tCopy = *task
		result = append(result, tCopy)
	}
	return result, nil
}
