package memory

import (
	"context"
	"errors"
	"hash/fnv"
	"sync"

	"github.com/gaz358/myprog/workmate/domain"
)

const shardCount = 16

type shard struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

type InMemoryRepo struct {
	shards [shardCount]shard
}

func NewInMemoryRepo() *InMemoryRepo {
	repo := &InMemoryRepo{}
	for i := range repo.shards {
		repo.shards[i].tasks = make(map[string]*domain.Task)
	}
	return repo
}

// Хэш-функция для распределения ID по шардам
func fnvHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// Возвращает нужный shard по id задачи
func (r *InMemoryRepo) getShard(id string) *shard {
	idx := fnvHash(id) % uint32(shardCount)
	return &r.shards[idx]
}

func (r *InMemoryRepo) Create(ctx context.Context, task *domain.Task) error {
	sh := r.getShard(task.ID)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if _, exists := sh.tasks[task.ID]; exists {
		return errors.New("task already exists")
	}
	tCopy := *task
	sh.tasks[task.ID] = &tCopy
	return nil
}

func (r *InMemoryRepo) Update(ctx context.Context, task *domain.Task) error {
	sh := r.getShard(task.ID)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if _, exists := sh.tasks[task.ID]; !exists {
		return domain.ErrNotFound
	}
	tCopy := *task
	sh.tasks[task.ID] = &tCopy
	return nil
}

func (r *InMemoryRepo) Delete(ctx context.Context, id string) error {
	sh := r.getShard(id)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if _, ok := sh.tasks[id]; !ok {
		return domain.ErrNotFound
	}
	delete(sh.tasks, id)
	return nil
}

func (r *InMemoryRepo) Get(ctx context.Context, id string) (*domain.Task, error) {
	sh := r.getShard(id)
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	t, ok := sh.tasks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	tCopy := *t
	return &tCopy, nil
}

func (r *InMemoryRepo) List(ctx context.Context) ([]*domain.Task, error) {
	result := make([]*domain.Task, 0)
	for i := 0; i < shardCount; i++ {
		sh := &r.shards[i]
		sh.mu.RLock()
		for _, task := range sh.tasks {
			tCopy := *task
			result = append(result, &tCopy)
		}
		sh.mu.RUnlock()
	}
	return result, nil
}
