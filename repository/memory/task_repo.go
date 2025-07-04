package memory

import (
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
		return domen.ErrNotFound
	}
	r.tasks[t.ID] = t
	return nil
}

func (r *InMemoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return domen.ErrNotFound
	}
	delete(r.tasks, id)
	return nil
}

func (r *InMemoryRepo) Get(id string) (*domen.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, domen.ErrNotFound
	}
	tCopy := *t // поверхностная копия!
	return &tCopy, nil
}

func (r *InMemoryRepo) List() ([]*domen.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*domen.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		tCopy := *t // поверхностная копия!
		tasks = append(tasks, &tCopy)
	}
	return tasks, nil
}
