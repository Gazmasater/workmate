package memory

import (
	"errors"
	"sync"

	"github.com/gaz358/myprog/workmate/domain"
)

type InMemoryRepo struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{tasks: make(map[string]*domain.Task)}
}

func (r *InMemoryRepo) Create(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; exists {
		return errors.New("task already exists")
	}
	tCopy := *task
	r.tasks[task.ID] = &tCopy
	return nil
}

func (r *InMemoryRepo) Update(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; !exists {
		return domain.ErrNotFound
	}
	tCopy := *task
	r.tasks[task.ID] = &tCopy
	return nil
}

func (r *InMemoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.tasks, id)
	return nil
}

func (r *InMemoryRepo) Get(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	tCopy := *t
	return &tCopy, nil
}

func (r *InMemoryRepo) List() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tCopy := *task
		result = append(result, &tCopy)
	}
	return result, nil
}
