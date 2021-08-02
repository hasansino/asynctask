[![Build Status](https://travis-ci.com/hasansino/asynctask.svg?branch=master)](https://travis-ci.com/hasansino/asynctask)
[![Go Report Card](https://goreportcard.com/badge/github.com/hasansino/asynctask)](https://goreportcard.com/report/github.com/hasansino/asynctask)

# asynctask
Simple async task runner.  
Executes provided functions (func() error) in asynchronous manner.  
Runner can be started in two fashions - blocking (Run()) and non-blocking (RunAsync()).

# Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/hasansino/asynctask"
)

func main() {
	runner := asynctask.New()

	runner.Add("task1", func(ctx context.Context) error {
		for i := 0; i < 1000; i++ {
			// do some job
		}
		return nil // return error
	})

	// Run() will block until all operations fill finish
	results := runner.Run()
	for _, r := range results {
		fmt.Printf("%s | %f | %v \n", r.Name(), r.Time().Seconds(), r.Error())
	}

	// RunAsync() will not block and will return channel to read
	resultCh := runner.RunAsync()
	for r := range resultCh {
		fmt.Printf("%s | %f | %v \n", r.Name(), r.Time().Seconds(), r.Error())
	}
}
```