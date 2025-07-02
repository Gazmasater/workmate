package usecase

import (
	"time"

	"github.com/gaz358/myprog/workmate/domen"
	"github.com/google/uuid"
)

type TaskUseCase struct {
	repo domen.TaskRepository
}

func NewTaskUseCase(repo domen.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
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
	time.Sleep(3 * time.Minute)
	task.Status = domen.StatusCompleted
	task.EndedAt = time.Now()
	task.Result = "OK"
	_ = uc.repo.Update(task)
}

func (uc *TaskUseCase) GetTask(id string) (*domen.Task, error) {
	return uc.repo.Get(id)
}

func (uc *TaskUseCase) DeleteTask(id string) error {
	return uc.repo.Delete(id)
}
