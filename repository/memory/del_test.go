package memory

import (
	"context"
	"testing"
	"time"

	"github.com/gaz358/myprog/workmate/domain"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepo_Delete(t *testing.T) {

	ctx := context.Background() // Можно объявить в начале теста, если его ещё нет

	repo := NewInMemoryRepo()

	task := &domain.Task{
		ID:        "task-to-delete",
		CreatedAt: time.Now(),
		Status:    domain.StatusPending,
	}

	// Создаем задачу
	err := repo.Create(ctx, task)
	assert.NoError(t, err, "ошибка при создании задачи")

	// Удаляем задачу
	err = repo.Delete(ctx, task.ID)
	assert.NoError(t, err, "ошибка при удалении задачи")

	// Проверяем, что задача действительно удалена
	_, err = repo.Get(ctx, task.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound, "ожидалась ошибка ErrNotFound после удаления")

	// Попытка удалить несуществующую задачу
	err = repo.Delete(ctx, "non-existent-id")
	assert.ErrorIs(t, err, domain.ErrNotFound, "ожидалась ошибка ErrNotFound при удалении несуществующей задачи")
}
