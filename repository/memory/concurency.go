package memory

import (
	"fmt"
	"sync"
	"testing"

	"github.com/gaz358/myprog/workmate/domen"
)

func TestInMemoryRepo_Concurrency(t *testing.T) {
	repo := NewInMemoryRepo()
	const n = 100
	var wg sync.WaitGroup

	// Параллельное создание задач
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tid := fmt.Sprintf("task-%d", id)
			task := &domen.Task{
				ID:     tid,
				Status: domen.StatusPending,
			}
			if err := repo.Create(task); err != nil {
				t.Errorf("create err: %v", err)
			}
		}(i)
	}
	wg.Wait()

	// Проверяем, что все задачи создались
	tasks, err := repo.List()
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(tasks) != n {
		t.Fatalf("want %d, got %d", n, len(tasks))
	}

	// Параллельное обновление задач
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tid := fmt.Sprintf("task-%d", id)
			task := &domen.Task{
				ID:     tid,
				Status: domen.StatusCompleted,
			}
			if err := repo.Update(task); err != nil {
				t.Errorf("update err: %v", err)
			}
		}(i)
	}
	wg.Wait()

	// Параллельное удаление задач
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tid := fmt.Sprintf("task-%d", id)
			if err := repo.Delete(tid); err != nil {
				t.Errorf("delete err: %v", err)
			}
		}(i)
	}
	wg.Wait()

	// После удаления ничего не должно остаться
	tasks, err = repo.List()
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("want 0, got %d", len(tasks))
	}
}
