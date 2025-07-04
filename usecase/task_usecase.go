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
	go uc.run(task.ID)
	// Возвращаем только копию задачи из репозитория, чтобы не светить указатель, который мутирует run
	return uc.repo.Get(task.ID)
}

func (uc *TaskUseCase) run(id string) {
	// Получаем копию задачи из репозитория
	task, err := uc.repo.Get(id)
	if err != nil {
		return
	}
	// Меняем статус и время старта
	task.Status = domen.StatusRunning
	task.StartedAt = time.Now()
	uc.repo.Update(task)

	time.Sleep(uc.duration)

	// Обновляем задачу после выполнения
	task.Status = domen.StatusCompleted
	task.EndedAt = time.Now()
	task.Duration = task.EndedAt.Sub(task.StartedAt).String()
	task.Result = "OK"
	uc.repo.Update(task)
}

func (uc *TaskUseCase) GetTask(id string) (*domen.Task, error) {
	return uc.repo.Get(id) // repo.Get должен возвращать копию!
}

func (uc *TaskUseCase) DeleteTask(id string) error {
	return uc.repo.Delete(id)
}

func (uc *TaskUseCase) ListTasks() ([]*domen.Task, error) {
	return uc.repo.List() // repo.List должен возвращать только копии!
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
	task.Result = "Canceled"
	return uc.repo.Update(task)
}
