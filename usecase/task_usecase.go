package usecase

import (
	"time"

	"github.com/gaz358/myprog/workmate/domen"
	"github.com/google/uuid"
)

type TaskUseCase struct {
	repo     domen.TaskRepository
	duration time.Duration
}

func NewTaskUseCase(repo domen.TaskRepository, duration time.Duration) *TaskUseCase {
	return &TaskUseCase{
		repo:     repo,
		duration: duration,
	}
}

func (uc *TaskUseCase) CreateTask() (*domen.Task, error) {
	task := &domen.Task{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		Status:    domen.StatusPending,
	}
	if err := uc.repo.Create(task); err != nil {
		return nil, err
	}
	go uc.run(task)
	return task, nil
}

func (uc *TaskUseCase) run(task *domen.Task) {
	task.Status = domen.StatusRunning
	task.StartedAt = time.Now()

	//time.Sleep(uc.duration)

	task.Status = domen.StatusCompleted
	task.EndedAt = time.Now()
	task.Duration = task.EndedAt.Sub(task.StartedAt).String()
	task.Result = "OK"

	_ = uc.repo.Update(task)
}

func (uc *TaskUseCase) GetTask(id string) (*domen.Task, error) {
	return uc.repo.Get(id)
}

func (uc *TaskUseCase) DeleteTask(id string) error {
	return uc.repo.Delete(id)
}

func (uc *TaskUseCase) ListTasks() ([]*domen.Task, error) {
	return uc.repo.List()
}

func (uc *TaskUseCase) CancelTask(id string) error {
	task, err := uc.repo.Get(id)
	if err != nil {
		return err
	}
	if task.Status == domen.StatusCompleted || task.Status == domen.StatusFailed || task.Status == domen.StatusCancelled {
		return nil
	}
	task.Status = domen.StatusCancelled
	task.Result = "Cancelled"
	return uc.repo.Update(task)
}
