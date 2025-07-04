package memory

import (
	"errors"
	"sync"

	"github.com/gaz358/myprog/workmate/domen"
)

type InMemoryRepo struct {
	mu    sync.RWMutex
	tasks map[string]*domen.Task
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{tasks: make(map[string]*domen.Task)}
}

func (r *InMemoryRepo) Create(task *domen.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; exists {
		return errors.New("task already exists")
	}
	// Копируем задачу
	tCopy := *task
	r.tasks[task.ID] = &tCopy
	return nil
}

func (r *InMemoryRepo) Update(task *domen.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; !exists {
		return errors.New("not found")
	}
	tCopy := *task
	r.tasks[task.ID] = &tCopy
	return nil
}

func (r *InMemoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return domen.ErrNotFound // <--- и тут тоже
	}
	delete(r.tasks, id)
	return nil
}

func (r *InMemoryRepo) Get(id string) (*domen.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, domen.ErrNotFound // <--- используем общую ошибку!
	}
	tCopy := *t
	return &tCopy, nil
}

func (r *InMemoryRepo) List() ([]*domen.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domen.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tCopy := *task
		result = append(result, &tCopy)
	}
	return result, nil
}
