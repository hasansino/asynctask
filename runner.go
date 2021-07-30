package asynctask

import (
	"sync"
	"time"
)

// AsyncRunner is async executor of tasks
type AsyncRunner struct {
	tasks []task
}

// New creates new empty async runner
func New() *AsyncRunner {
	return &AsyncRunner{tasks: make([]task, 0)}
}

// Add function to execute async
func (l *AsyncRunner) Add(name string, fn taskFn) {
	l.tasks = append(l.tasks, task{
		name: name,
		fn:   fn,
	})
}

// Run all async loaders and collect results
// Returns slice of errors collected for all executed tasks
func (l *AsyncRunner) Run() []TaskResult {
	ch := l.runParallel()

	results := make([]TaskResult, 0, len(l.tasks))

	for res := range ch {
		results = append(results, res)
	}

	return results
}

// RunAsync executes tasks in non-blocking fashion
// Returns channel to read results from
func (l *AsyncRunner) RunAsync() <-chan TaskResult {
	return l.runParallel()
}

// runParallel executes tasks parallel and return channel
// where results will be sent once tasks are finished.
func (l *AsyncRunner) runParallel() <-chan TaskResult {
	wg := new(sync.WaitGroup)
	wg.Add(len(l.tasks))

	ch := make(chan TaskResult, len(l.tasks))

	for i := range l.tasks {
		go func(ch chan TaskResult, l task) {
			defer wg.Done()
			var err error
			timeStart := time.Now()
			err = l.fn()
			execTime := time.Since(timeStart)
			ch <- TaskResult{
				name:     l.name,
				execTime: execTime,
				err:      err,
			}
		}(ch, l.tasks[i])
	}

	// wait till all tasks finish
	wg.Wait()
	// close channel, since we are not expecting any writes
	close(ch)

	return ch
}
