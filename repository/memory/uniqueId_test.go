package memory

import (
	"testing"
	"time"

	"github.com/gaz358/myprog/workmate/usecase"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask_UniqueIDs(t *testing.T) {
	repo := NewInMemoryRepo()
	uc := usecase.NewTaskUseCase(repo, 1*time.Second)

	task1, err1 := uc.CreateTask()
	task2, err2 := uc.CreateTask()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, task1.ID, task2.ID, "ожидаются уникальные ID")
}
