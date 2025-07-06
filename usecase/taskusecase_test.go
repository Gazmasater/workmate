package usecase

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gaz358/myprog/workmate/domain"
	"github.com/stretchr/testify/assert"
)

type fakeRepo struct {
	mu        sync.Mutex
	tasks     map[string]*domain.Task
	createErr error
	updateErr error
	deleteErr error
	getErr    error
	listErr   error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{tasks: make(map[string]*domain.Task)}
}

func (f *fakeRepo) Create(ctx context.Context, t *domain.Task) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	copy := *t
	f.tasks[t.ID] = &copy
	return nil
}
func (f *fakeRepo) Update(ctx context.Context, t *domain.Task) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	copy := *t
	f.tasks[t.ID] = &copy
	return nil
}
func (f *fakeRepo) Get(ctx context.Context, id string) (*domain.Task, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	t, ok := f.tasks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copy := *t
	return &copy, nil
}

func (f *fakeRepo) Delete(ctx context.Context, id string) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.tasks, id)
	return nil
}

func (f *fakeRepo) List(ctx context.Context) ([]*domain.Task, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]*domain.Task, 0, len(f.tasks))
	for _, t := range f.tasks {
		copy := *t
		out = append(out, &copy)
	}
	return out, nil
}

func TestTaskUseCase_CreateTask_Success(t *testing.T) {
	repo := newFakeRepo()
	uc := NewTaskUseCase(repo, 10*time.Millisecond)

	task, err := uc.CreateTask()
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, domain.StatusPending, task.Status)
	assert.NotEmpty(t, task.ID)

	// Проверяем, что задача появилась в fakeRepo
	repo.mu.Lock()
	_, ok := repo.tasks[task.ID]
	repo.mu.Unlock()
	assert.True(t, ok)
}

func TestTaskUseCase_CreateTask_RepoError(t *testing.T) {
	repo := newFakeRepo()
	repo.createErr = assert.AnError
	uc := NewTaskUseCase(repo, 10*time.Millisecond)

	task, err := uc.CreateTask()
	assert.Error(t, err)
	assert.Nil(t, task)
}

func TestTaskUseCase_GetTask(t *testing.T) {
	repo := newFakeRepo()
	task := &domain.Task{ID: "tid"}
	repo.tasks[task.ID] = task

	uc := NewTaskUseCase(repo, 10*time.Millisecond)
	got, err := uc.GetTask("tid")
	assert.NoError(t, err)
	assert.Equal(t, task.ID, got.ID)

	_, err = uc.GetTask("notfound")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTaskUseCase_ListTasks(t *testing.T) {
	repo := newFakeRepo()
	repo.tasks["1"] = &domain.Task{ID: "1"}
	repo.tasks["2"] = &domain.Task{ID: "2"}
	uc := NewTaskUseCase(repo, 10*time.Millisecond)

	tasks, err := uc.ListTasks()
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
}

func TestTaskUseCase_DeleteTask(t *testing.T) {
	repo := newFakeRepo()
	task := &domain.Task{ID: "tid"}
	repo.tasks[task.ID] = task
	uc := NewTaskUseCase(repo, 10*time.Millisecond)

	// Создадим cancelMap вручную (для теста отмены)
	uc.cancelMap[task.ID] = func() {}

	err := uc.DeleteTask(task.ID)
	assert.NoError(t, err)

	repo.mu.Lock()
	_, ok := repo.tasks[task.ID]
	repo.mu.Unlock()
	assert.False(t, ok)
}

func TestTaskUseCase_CancelTask(t *testing.T) {
	repo := newFakeRepo()
	task := &domain.Task{ID: "tid"}
	repo.tasks[task.ID] = task
	uc := NewTaskUseCase(repo, 10*time.Millisecond)

	// Проверяем успешную отмену
	called := false
	uc.cancelMap[task.ID] = func() { called = true }
	err := uc.CancelTask(task.ID)
	assert.NoError(t, err)
	assert.True(t, called)

	// Проверяем отмену несуществующей задачи
	err = uc.CancelTask("notfound")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
