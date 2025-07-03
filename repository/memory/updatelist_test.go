package memory

import (
	"testing"
	"time"

	"github.com/gaz358/myprog/workmate/domen"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepo_Update(t *testing.T) {
	repo := NewInMemoryRepo()

	// Создание и добавление задачи
	task := &domen.Task{
		ID:        "task-1",
		CreatedAt: time.Now(),
		Status:    domen.StatusPending,
	}
	err := repo.Create(task)
	assert.NoError(t, err, "ошибка при создании задачи")

	// Обновление задачи
	task.Status = domen.StatusCompleted
	task.Result = "done"
	err = repo.Update(task)
	assert.NoError(t, err, "ошибка при обновлении задачи")

	updated, err := repo.Get(task.ID)
	assert.NoError(t, err)
	assert.Equal(t, domen.StatusCompleted, updated.Status)
	assert.Equal(t, "done", updated.Result)

	// Попытка обновить несуществующую задачу
	nonexistent := &domen.Task{
		ID:     "nonexistent",
		Status: domen.StatusFailed,
	}
	err = repo.Update(nonexistent)
	assert.ErrorIs(t, err, domen.ErrNotFound)
}

func TestInMemoryRepo_List(t *testing.T) {
	repo := NewInMemoryRepo()

	// Пустой список
	tasks, err := repo.List()
	assert.NoError(t, err)
	assert.Empty(t, tasks, "список должен быть пуст при отсутствии задач")

	// Добавим несколько задач
	task1 := &domen.Task{ID: "id1", CreatedAt: time.Now(), Status: domen.StatusPending}
	task2 := &domen.Task{ID: "id2", CreatedAt: time.Now(), Status: domen.StatusCompleted}

	_ = repo.Create(task1)
	_ = repo.Create(task2)

	tasks, err = repo.List()
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)

	ids := map[string]bool{tasks[0].ID: true, tasks[1].ID: true}
	assert.True(t, ids["id1"])
	assert.True(t, ids["id2"])
}
