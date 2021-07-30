package asynctask

import "time"

// taskFn is loaders that are accepted by asyncLoader
type taskFn func() error

// task is a single function to load some arbitrary data
type task struct {
	name string
	fn   taskFn
}

// TaskResult is async task result which is passed to channel
type TaskResult struct {
	name     string        // name of task
	execTime time.Duration // execution duration
	err      error         // resulting error
}

// Name of task
func (ts *TaskResult) Name() string {
	return ts.name
}

// Time is duration of task execution
func (ts *TaskResult) Time() time.Duration {
	return ts.execTime
}

// Error result of task
func (ts *TaskResult) Error() error {
	return ts.err
}
