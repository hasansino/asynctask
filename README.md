[![Go Report Card](https://goreportcard.com/badge/github.com/hasansino/asynctask)](https://goreportcard.com/report/github.com/hasansino/asynctask)
[![Build Status](https://travis-ci.com/hasansino/asynctask.svg?branch=master)](https://travis-ci.com/hasansino/asynctask)

# asynctask

asynctask executes tasks (func(context.Context) error) asynchronously with optional context timeout.

## Installation

```bash
~ $ go get -u github.com/hasansino/asynctask
```

## Example

```go
package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/hasansino/asynctask"
)

func main() {
	runner := asynctask.New()
	runner.Add("example", func(ctx context.Context) error {
		for i := 0; i < 1000; i++ {
			math.Log(rand.Float64()) // do some job
		}
		return nil
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

## Example with timeout

```go
package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hasansino/asynctask"
)

func main() {
	runner := asynctask.New()
	runner.Add("example", func(ctx context.Context) error {
		for {
			switch err := ctx.Err(); {
			case err != nil:
				return err // context is canceled or deadline reached
			default:
				math.Log(rand.Float64()) // do some job
			}
		}
	})

	// Set timeout of 1 second
	runner.SetTimeout(time.Second)

	// Run() will block until all operations fill finish
	results := runner.Run()
	for _, r := range results {
		fmt.Printf("%s | %f | %v \n", r.Name(), r.Time().Seconds(), r.Error())
	}
}
```

## Example with middleware

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/hasansino/asynctask"
)

func main() {
	runner := asynctask.New()

	runner.Use(func(next asynctask.TaskFn) asynctask.TaskFn {
		return func(ctx context.Context) error {
			log.Println("Before `example` task")
			err := next(ctx)
			log.Println("After `example` task")
			return err
		}
	})

	runner.Add("example", func(ctx context.Context) error {
		for i := 0; i < 10000; i++ {
			math.Log(rand.Float64()) // do some job
		}
		return nil
	})

	// Run() will block until all operations fill finish
	results := runner.Run()
	for _, r := range results {
		fmt.Printf("%s | %f | %v \n", r.Name(), r.Time().Seconds(), r.Error())
	}
}
```