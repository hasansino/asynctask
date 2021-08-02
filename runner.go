package asynctask

import (
	"context"
	"sync"
	"time"
)

// AsyncRunner is async executor of tasks
type AsyncRunner struct {
	timeout    time.Duration
	tasks      []Task
	middleware []MiddlewareFn
}

// New creates new empty async runner
func New() *AsyncRunner {
	return &AsyncRunner{
		tasks:      make([]Task, 0),
		middleware: make([]MiddlewareFn, 0),
	}
}

// Add Task to execution list
func (r *AsyncRunner) Add(name string, fn TaskFn) {
	r.tasks = append(r.tasks, Task{
		name: name,
		fn:   fn,
	})
}

// Use adds middleware to the chain
func (r *AsyncRunner) Use(middleware ...MiddlewareFn) {
	for _, m := range middleware {
		r.middleware = append(r.middleware, m)
	}
}

// SetTimeout for all tasks, unlimited if zero
// When deadline is reached CancelFn is called on context
// and task provider is responsible for handling it properly
func (r *AsyncRunner) SetTimeout(t time.Duration) {
	r.timeout = t
}

// Reset runner to initial empty state
func (r *AsyncRunner) Reset() {
	r.timeout = 0
	r.tasks = make([]Task, 0)
	r.middleware = make([]MiddlewareFn, 0)
}

// Run all async loaders and collect results
// Returns slice of errors collected for all executed tasks
func (r *AsyncRunner) Run() []TaskResult {
	ch := r.runParallel()

	results := make([]TaskResult, 0, len(r.tasks))

	for res := range ch {
		results = append(results, res)
	}

	return results
}

// RunAsync executes tasks in non-blocking fashion
// Returns channel to read results from
func (r *AsyncRunner) RunAsync() <-chan TaskResult {
	return r.runParallel()
}

// runParallel executes tasks parallel and return channel
// where results will be sent once tasks are finished.
func (r *AsyncRunner) runParallel() <-chan TaskResult {
	wg := new(sync.WaitGroup)
	ch := make(chan TaskResult, len(r.tasks))

	ctx := context.Background()
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	wg.Add(len(r.tasks))
	for i := range r.tasks {
		go r.subTask(ctx, wg, ch, &r.tasks[i])
	}

	// wait till all tasks will finish
	wg.Wait()
	// close channel, since we are not expecting any writes
	close(ch)

	return ch
}

// subTask is background task routine
func (r *AsyncRunner) subTask(
	ctx context.Context,
	wg *sync.WaitGroup, ch chan TaskResult, t *Task,
) {
	defer wg.Done()

	var err error
	ts := time.Now()

	fn := t.fn
	for i := 0; i < len(r.middleware); i++ {
		fn = r.middleware[i](fn)
	}

	err = fn(ctx)

	et := time.Since(ts)

	ch <- TaskResult{
		name: t.name,
		time: et,
		err:  err,
	}
}
