package thermocline

import (
	"errors"
	"time"

	"github.com/satori/go.uuid"
)

// MaxRetries is the number of times a task will be attempted before it's marked
// "failed"
var MaxRetries = 3

// DefaultPriority is the default priority of tasks
var DefaultPriority = 1

// A Task contains information sent from the queue-er to the worker
type Task struct {
	ID      string `json:"id"`
	Retries int    `json:"retries"`

	// CreatedAt
	CreatedAt time.Time `json:"created_at"`
	// WorkAt is the earliest this task can be worked
	WorkAt   time.Time `json:"work_at"`
	Priority int       `json:"priority"`

	Info interface{} `json:"info"`
}

// NewTask returns a normal old task
func NewTask(info interface{}) (*Task, error) {
	now := time.Now()
	return &Task{
		ID:        uuid.NewV4().String(),
		Retries:   0,
		CreatedAt: now,
		WorkAt:    now,
		Priority:  DefaultPriority,
		Info:      info,
	}, nil
}

// NewDelayedTask returns a task that will be worked after a certain time has
// elapsed
func NewDelayedTask(info interface{}, delay time.Duration) (*Task, error) {
	if delay.Nanoseconds() < 0 {
		return nil, errors.New("cannot create a task delayed into the past")
	}

	now := time.Now()
	return &Task{
		ID:        uuid.NewV4().String(),
		Retries:   0,
		CreatedAt: now,
		WorkAt:    now.Add(delay),
		Priority:  DefaultPriority,
		Info:      info,
	}, nil
}
