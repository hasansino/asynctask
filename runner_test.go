package asynctask

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _testWorkFn = func() {
	for i := float64(0); i < 1000; i++ {
		_ = math.Log(i)
	}
}

func TestRunner_Success(t *testing.T) {
	tasks := []task{
		{"test", func() error { _testWorkFn(); return nil }},
		{"test1", func() error { _testWorkFn(); return nil }},
		{"test2", func() error { _testWorkFn(); return nil }},
		{"test3", func() error { _testWorkFn(); return nil }},
		{"test4", func() error { _testWorkFn(); return nil }},
		{"test5", func() error { _testWorkFn(); return nil }},
		{"test6", func() error { _testWorkFn(); return nil }},
		{"test7", func() error { _testWorkFn(); return nil }},
	}

	runner := New()

	for _, t := range tasks {
		runner.Add(t.name, t.fn)
	}

	results := runner.Run()
	assert.Len(t, results, len(tasks))
	assert.False(t, _testResultsHaveError(results))
}

func TestRunner_Error(t *testing.T) {
	tasks := []task{
		{"test", func() error { _testWorkFn(); return fmt.Errorf("fail") }},
	}

	runner := New()

	for _, t := range tasks {
		runner.Add(t.name, t.fn)
	}

	results := runner.Run()
	assert.Len(t, results, len(tasks))
	assert.True(t, _testResultsHaveError(results))
}

func TestRunnerAsync(t *testing.T) {
	tasks := []task{
		{"test", func() error { _testWorkFn(); return fmt.Errorf("fail") }},
	}

	runner := New()

	for _, t := range tasks {
		runner.Add(t.name, t.fn)
	}

	ch := runner.RunAsync()

	results := make([]TaskResult, 0)
	for res := range ch {
		results = append(results, res)
	}

	assert.Len(t, results, len(tasks))
	assert.True(t, _testResultsHaveError(results))
}

func _testResultsHaveError(results []TaskResult) bool {
	for _, r := range results {
		if r.err != nil {
			return true
		}
	}
	return false
}
