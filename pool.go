package asynctask

import "sync"

var runnerPool = sync.Pool{
	New: func() interface{} {
		return New()
	},
}

// AcquireRunner from sync.Pool
func AcquireRunner() *AsyncRunner {
	return runnerPool.Get().(*AsyncRunner)
}

// ReleaseRunner back to sync.Pool and reset its state
func ReleaseRunner(r *AsyncRunner) {
	r.Reset()
	runnerPool.Put(r)
}
