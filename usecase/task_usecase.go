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

// NewTaskUseCase создает новый экземпляр TaskUseCase.
func NewTaskUseCase(repo domen.TaskRepository, duration time.Duration) *TaskUseCase {
	return &TaskUseCase{
		repo:     repo,
		duration: duration,
	}
}

// CreateTask создает новую задачу, сохраняет её в репозитории и запускает выполнение в отдельной горутине.
// Возвращает копию задачи из репозитория.
func (uc *TaskUseCase) CreateTask() (*domen.Task, error) {
	task := &domen.Task{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		Status:    domen.StatusPending,
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

// run выполняет задачу с заданным идентификатором.
// Все изменения статуса задачи сохраняются в репозитории с проверкой ошибок.
func (uc *TaskUseCase) run(id string) {
	task, err := uc.repo.Get(id)
	if err != nil {
		// Не удалось получить задачу — выходим.
		return
	}

	task.Status = domen.StatusRunning
	task.StartedAt = time.Now()
	err = uc.repo.Update(task)
	if err != nil {
		// Не удалось обновить статус задачи — выходим.
		return
	}

	time.Sleep(uc.duration)

	task.Status = domen.StatusCompleted
	task.EndedAt = time.Now()
	task.Duration = task.EndedAt.Sub(task.StartedAt).String()
	task.Result = "OK"
	err = uc.repo.Update(task)
	if err != nil {
		// Не удалось обновить задачу после выполнения — выходим.
		return
	}
}

// GetTask возвращает копию задачи по идентификатору.
func (uc *TaskUseCase) GetTask(id string) (*domen.Task, error) {
	task, err := uc.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// DeleteTask удаляет задачу по идентификатору.
func (uc *TaskUseCase) DeleteTask(id string) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// ListTasks возвращает список копий всех задач.
func (uc *TaskUseCase) ListTasks() ([]*domen.Task, error) {
	tasks, err := uc.repo.List()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// CancelTask отменяет задачу по идентификатору, если она не завершена.
func (uc *TaskUseCase) CancelTask(id string) error {
	task, err := uc.repo.Get(id)
	if err != nil {
		return err
	}
	if task.Status == domen.StatusCompleted ||
		task.Status == domen.StatusFailed ||
		task.Status == domen.StatusCancelled {
		return nil
	}
	task.Status = domen.StatusCancelled
	task.Result = "Canceled"
	err = uc.repo.Update(task)
	if err != nil {
		return err
	}
	return nil
}
