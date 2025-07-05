package memory

import (
	"context"
	"testing"
	"time"

	"github.com/gaz358/myprog/workmate/domain"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepo_CreateAndGet_ExactMatch(t *testing.T) {
	ctx := context.Background() // Можно объявить в начале теста, если его ещё нет

	repo := NewInMemoryRepo()

	expectedTask := &domain.Task{
		ID:        "task-abc123",
		CreatedAt: time.Now().Truncate(time.Second),
		StartedAt: time.Now().Add(1 * time.Second).Truncate(time.Second),
		EndedAt:   time.Now().Add(5 * time.Second).Truncate(time.Second),
		Duration:  "4s",
		Status:    domain.StatusCompleted,
		Result:    "OK",
	}

	err := repo.Create(ctx, expectedTask)
	assert.NoError(t, err, "ошибка при создании задачи")

	got, err := repo.Get(ctx, expectedTask.ID)
	assert.NoError(t, err, "ошибка при получении задачи")

	assert.Equal(t, expectedTask.ID, got.ID)
	assert.Equal(t, expectedTask.CreatedAt, got.CreatedAt)
	assert.Equal(t, expectedTask.StartedAt, got.StartedAt)
	assert.Equal(t, expectedTask.EndedAt, got.EndedAt)
	assert.Equal(t, expectedTask.Duration, got.Duration)
	assert.Equal(t, expectedTask.Status, got.Status)
	assert.Equal(t, expectedTask.Result, got.Result)

	// Проверка получения несуществующей задачи
	nonExistentID := "task-nonexistent"
	got, err = repo.Get(ctx, nonExistentID)
	assert.Nil(t, got, "ожидается nil при получении несуществующей задачи")
	assert.ErrorIs(t, err, domain.ErrNotFound, "ожидалась ошибка ErrNotFound")
}
