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

func (r *InMemoryRepo) Create(t *domen.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[t.ID] = t
	return nil
}

func (r *InMemoryRepo) Update(t *domen.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[t.ID]; !ok {
		return errors.New("not found")
	}
	r.tasks[t.ID] = t
	return nil
}

func (r *InMemoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tasks, id)
	return nil
}

func (r *InMemoryRepo) Get(id string) (*domen.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}
