package domain

type TaskRepository interface {
	Create(*Task) error
	Update(*Task) error
	Delete(id string) error
	Get(id string) (*Task, error)
}
