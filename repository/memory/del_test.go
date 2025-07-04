package memory

import (
	"testing"
	"time"

	"github.com/gaz358/myprog/workmate/domen"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepo_Delete(t *testing.T) {
	repo := NewInMemoryRepo()

	task := &domen.Task{
		ID:        "task-to-delete",
		CreatedAt: time.Now(),
		Status:    domen.StatusPending,
	}

	// Создаем задачу
	err := repo.Create(task)
	assert.NoError(t, err, "ошибка при создании задачи")

	// Удаляем задачу
	err = repo.Delete(task.ID)
	assert.NoError(t, err, "ошибка при удалении задачи")

	// Проверяем, что задача действительно удалена
	_, err = repo.Get(task.ID)
	assert.ErrorIs(t, err, domen.ErrNotFound, "ожидалась ошибка ErrNotFound после удаления")

	// Попытка удалить несуществующую задачу
	err = repo.Delete("non-existent-id")
	assert.ErrorIs(t, err, domen.ErrNotFound, "ожидалась ошибка ErrNotFound при удалении несуществующей задачи")
}
