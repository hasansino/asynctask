package asynctask

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _testWorkFn = func() {
	for i := float64(0); i < 1000; i++ {
		_ = math.Log(i)
	}
}

func _testResultsHaveError(results []TaskResult) bool {
	for _, r := range results {
		if r.err != nil {
			return true
		}
	}
	return false
}

func TestRunner_Success(t *testing.T) {
	tasks := []Task{
		{"test1", func(ctx context.Context) error { _testWorkFn(); return nil }},
		{"test2", func(ctx context.Context) error { _testWorkFn(); return nil }},
		{"test3", func(ctx context.Context) error { _testWorkFn(); return nil }},
		{"test4", func(ctx context.Context) error { _testWorkFn(); return nil }},
		{"test5", func(ctx context.Context) error { _testWorkFn(); return nil }},
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
	tasks := []Task{
		{"test1", func(ctx context.Context) error { _testWorkFn(); return fmt.Errorf("fail") }},
		{"test2", func(ctx context.Context) error { _testWorkFn(); return fmt.Errorf("fail") }},
		{"test3", func(ctx context.Context) error { _testWorkFn(); return fmt.Errorf("fail") }},
		{"test4", func(ctx context.Context) error { _testWorkFn(); return fmt.Errorf("fail") }},
		{"test5", func(ctx context.Context) error { _testWorkFn(); return fmt.Errorf("fail") }},
	}

	runner := New()
	for _, t := range tasks {
		runner.Add(t.name, t.fn)
	}

	results := runner.Run()
	assert.Len(t, results, len(tasks))
	assert.True(t, _testResultsHaveError(results))
}

func TestRunnerAsync_Success(t *testing.T) {
	tasks := []Task{
		{"test1", func(ctx context.Context) error { _testWorkFn(); return nil }},
		{"test2", func(ctx context.Context) error { _testWorkFn(); return nil }},
		{"test3", func(ctx context.Context) error { _testWorkFn(); return nil }},
		{"test4", func(ctx context.Context) error { _testWorkFn(); return nil }},
		{"test5", func(ctx context.Context) error { _testWorkFn(); return nil }},
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
	assert.False(t, _testResultsHaveError(results))
}

func TestRunnerAsync_Fail(t *testing.T) {
	tasks := []Task{
		{"test1", func(ctx context.Context) error { _testWorkFn(); return fmt.Errorf("fail") }},
		{"test2", func(ctx context.Context) error { _testWorkFn(); return fmt.Errorf("fail") }},
		{"test3", func(ctx context.Context) error { _testWorkFn(); return fmt.Errorf("fail") }},
		{"test4", func(ctx context.Context) error { _testWorkFn(); return fmt.Errorf("fail") }},
		{"test5", func(ctx context.Context) error { _testWorkFn(); return fmt.Errorf("fail") }},
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

func TestMiddleware(t *testing.T) {
	tasks := []Task{
		{"test1", func(ctx context.Context) error {
			_testWorkFn()
			ctxValue := ctx.Value("test")
			if ctxValue == nil {
				t.Error("ctx value is nil")
			}
			if strCtxValue, ok := ctxValue.(string); !ok {
				t.Error("ctx value have invalid type")
			} else if strCtxValue != "123" {
				t.Error("ctx value have invalid value")
			}
			return nil
		}},
	}

	runner := New()
	for _, t := range tasks {
		runner.Add(t.name, t.fn)
	}

	runner.Use(
		func(next TaskFn) TaskFn {
			return func(ctx context.Context) error {
				prevValue := ctx.Value("test").(string)
				ctx = context.WithValue(ctx, "test", fmt.Sprintf("%s3", prevValue))
				return next(ctx)
			}
		},
		func(next TaskFn) TaskFn {
			return func(ctx context.Context) error {
				prevValue := ctx.Value("test").(string)
				ctx = context.WithValue(ctx, "test", fmt.Sprintf("%s2", prevValue))
				return next(ctx)
			}
		},
		func(next TaskFn) TaskFn {
			return func(ctx context.Context) error {
				ctx = context.WithValue(ctx, "test", "1")
				return next(ctx)
			}
		},
	)

	results := runner.Run()
	assert.Len(t, results, len(tasks))
	assert.False(t, _testResultsHaveError(results))
}

func TestTimeout(t *testing.T) {
	tasks := []Task{
		{"test1", func(ctx context.Context) error {
			for {

				switch err := ctx.Err(); {
				case err != nil:
					return err
				default:
					_testWorkFn()
				}
			}
		}},
	}

	runner := New()
	for _, t := range tasks {
		runner.Add(t.name, t.fn)
	}

	runner.SetTimeout(time.Millisecond)

	results := runner.Run()
	assert.Len(t, results, len(tasks))
	assert.Equal(t, context.DeadlineExceeded, results[0].Error())
}
