package domen

type TaskRepository interface {
	Create(task *Task) error
	Update(task *Task) error
	Get(id string) (*Task, error)
	List() ([]*Task, error)
	Delete(id string) error
}
