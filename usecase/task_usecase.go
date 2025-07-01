package usecase

import (
	"time"
	"workmate/domain"

	"github.com/google/uuid"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) CreateTask() (*domain.Task, error) {
	task := &domain.Task{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		Status:    domain.StatusPending,
	}
	if err := uc.repo.Create(task); err != nil {
		return nil, err
	}
	go uc.run(task)
	return task, nil
}

func (uc *TaskUseCase) run(task *domain.Task) {
	task.Status = domain.StatusRunning
	task.StartedAt = time.Now()
	time.Sleep(3 * time.Minute)
	task.Status = domain.StatusCompleted
	task.EndedAt = time.Now()
	task.Result = "OK"
	_ = uc.repo.Update(task)
}

func (uc *TaskUseCase) GetTask(id string) (*domain.Task, error) {
	return uc.repo.Get(id)
}

func (uc *TaskUseCase) DeleteTask(id string) error {
	return uc.repo.Delete(id)
}
