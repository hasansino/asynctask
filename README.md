[![Build Status](https://travis-ci.com/hasansino/asynctask.svg?branch=master)](https://travis-ci.com/hasansino/asynctask)
[![Go Report Card](https://goreportcard.com/badge/github.com/hasansino/asynctask)](https://goreportcard.com/report/github.com/hasansino/asynctask)

# asynctask
Simple async task runner.  
Executes provided functions (func() error) in asynchronous manner.  
Runner can be started in two fashions - blocking (Run()) and non-blocking (RunAsync()).

# Usage

```go
package main

import "github.com/hasansino/asynctask"

func main() {
	runner := asynctask.New()

	runner.Add("task1", func() error {
        // do some job
		return nil // return error
	})
	
	// Run() will block until all operations fill finish
	results := runner.Run()
	
	// RunAsync() will not block and will return channel to read
	resultCh := runner.RunAsync()
}
```