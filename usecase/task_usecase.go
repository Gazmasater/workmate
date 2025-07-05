package usecase

import (
	"time"

	"github.com/gaz358/myprog/workmate/domain"
	"github.com/google/uuid"
)

type TaskUseCase struct {
	repo     domain.TaskRepository
	duration time.Duration
}

func NewTaskUseCase(repo domain.TaskRepository, duration time.Duration) *TaskUseCase {
	return &TaskUseCase{
		repo:     repo,
		duration: duration,
	}
}

func (uc *TaskUseCase) CreateTask() (*domain.Task, error) {
	task := &domain.Task{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		Status:    domain.StatusPending,
	}
	err := uc.repo.Create(task)
	if err != nil {
		return nil, err
	}

	go uc.run(task.ID)

	createdTask, err := uc.repo.Get(task.ID)
	if err != nil {
		return nil, err
	}
	return createdTask, nil
}

func (uc *TaskUseCase) run(id string) {
	task, err := uc.repo.Get(id)
	if err != nil {
		return
	}

	task.Status = domain.StatusRunning
	task.StartedAt = time.Now()
	err = uc.repo.Update(task)
	if err != nil {
		return
	}

	time.Sleep(uc.duration)

	task.Status = domain.StatusCompleted
	task.EndedAt = time.Now()
	task.Duration = task.EndedAt.Sub(task.StartedAt).String()
	task.Result = "OK"
	err = uc.repo.Update(task)
	if err != nil {
		return
	}
}

func (uc *TaskUseCase) GetTask(id string) (*domain.Task, error) {
	task, err := uc.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (uc *TaskUseCase) DeleteTask(id string) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *TaskUseCase) ListTasks() ([]*domain.Task, error) {
	tasks, err := uc.repo.List()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (uc *TaskUseCase) CancelTask(id string) error {
	task, err := uc.repo.Get(id)
	if err != nil {
		return err
	}
	if task.Status == domain.StatusCompleted ||
		task.Status == domain.StatusFailed ||
		task.Status == domain.StatusCancelled {
		return nil
	}
	task.Status = domain.StatusCancelled
	task.Result = "Canceled"
	err = uc.repo.Update(task)
	if err != nil {
		return err
	}
	return nil
}
