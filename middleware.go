package asynctask

// MiddlewareFn defines a function to process middleware
// Any middleware function provided for runner will may be
// called by multiple goroutines and thereby should be thread-safe
type MiddlewareFn func(TaskFn) TaskFn
