package asynctask

import (
	"context"
	"time"
)

// TaskFn is loaders that are accepted by asyncLoader
type TaskFn func(context.Context) error

// Task is a single function to load some arbitrary data
type Task struct {
	name string
	fn   TaskFn
}

// TaskResult is async Task result which is passed to channel
type TaskResult struct {
	name string        // name of Task
	time time.Duration // execution duration
	err  error         // resulting error
}

// Name of Task
func (ts *TaskResult) Name() string {
	return ts.name
}

// Time is duration of Task execution
func (ts *TaskResult) Time() time.Duration {
	return ts.time
}

// Error result of Task
func (ts *TaskResult) Error() error {
	return ts.err
}
