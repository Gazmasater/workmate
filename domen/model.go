package domen

import "time"

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusRunning   Status = "RUNNING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

// swagger:model Task
type Task struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	StartedAt time.Time `json:"started_at,omitempty"`
	EndedAt   time.Time `json:"ended_at,omitempty"`

	// Duration of the task execution
	// example: 3m0s
	Duration string `json:"duration,omitempty"`

	Status Status `json:"status"`
	Result string `json:"result,omitempty"`
}

// swagger:model TaskListItem
type TaskListItem struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Duration string `json:"duration,omitempty"`
}
