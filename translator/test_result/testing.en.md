---
title: Understanding the Difference Between log.Fatal() and panic() in Golang
type: blog
date: 2025-02-11
comments: true
---

## "In such cases, I think it's better to use `panic()` than `log.Fatal()"
Recently, I received feedback saying that it’s better to use `panic()` than `log.Fatal()`.

Hmm? Isn’t `log.Fatal()` just a better way to log errors? I thought.

It’s somewhat embarrassing, but it was only recently that I became clearly aware of the difference between `log.Fatal()` and `panic()` in Golang, so I’d like to organize my thoughts on it this time.

## The Difference Between log.Fatal() and panic()
Both `log.Fatal()` and `panic()` are functions that terminate a program. Let’s explore their behavior through the code.

```go
package main

import (
    "log"
	"log/slog"
)

func RunWithFatal() {
	log.Fatal("This is a fatal error")
}

func main() {
	RunWithFatal()

	slog.Info("This is not executed")
}
```

When we run the code above, we will see the following result:

```shell
2025/02/11 20:02:31 This is a fatal error
```

Now, let’s look at a code that uses `panic()`.

```go
package main

import (
    "log/slog"
)

func RunWithPanic() {
    panic("This is a panic error")
}

func main() {
    RunWithPanic()

    slog.Info("This is not executed")
}
```

When we run this code, we will see the following result:

```shell
panic: This is a panic error

goroutine 1 [running]:
main.RunWithPanic(...)
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:8
main.main()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:12 +0x30
```

From the above code, we see that `log.Fatal()` prints an error and terminates the program. In contrast, while `panic()` also prints an error and terminates the program, it outputs the stack trace as well.

### Recovering Using recover()
When using `panic()`, the program terminates. However, by using `recover()`, we can recover without terminating the program.

Ideally, it would be best if there were no panic at all, but developers are human, and mistakes are bound to happen. For that reason, in API servers and similar environments, middlewares are often created to recover from `panic()`, preventing situations where the server unexpectedly crashes.

To clearly see the difference, let’s first try to recover from `log.Fatal()`:
```go
package main

import (
    "log"
    "log/slog"
)

func RunWithFatal() {
    log.Fatal("This is a fatal error")
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            slog.Info("Recovered from", "error", r)
        }
    }()

    RunWithFatal()

    slog.Info("This is not executed")
}
```

When we run the above code, we will see the following result:

```shell
2025/02/11 20:07:49 This is a fatal error
```

It has not recovered. Now, let’s see if we can recover from `panic()`:

```go
package main

import (
	"log/slog"
	"runtime/debug"
)

func RunWithPanic() {
	panic("This is a panic error")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Info("Recovered from", "error", r)
			debug.PrintStack()
		}
	}()

	RunWithPanic()
}
```

When this code is executed, we receive the following output:

```shell
2025/02/11 20:09:51 INFO Recovered from error="This is a panic error"
goroutine 1 [running]:
runtime/debug.Stack()
	/opt/homebrew/opt/go/libexec/src/runtime/debug/stack.go:26 +0x64
runtime/debug.PrintStack()
	/opt/homebrew/opt/go/libexec/src/runtime/debug/stack.go:18 +0x1c
main.main.func1()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:16 +0x8c
panic({0x1004d6560?, 0x1004f4190?})
	/opt/homebrew/opt/go/libexec/src/runtime/panic.go:785 +0x124
main.RunWithPanic(...)
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:9
main.main()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:20 +0x4c
```

Although it is not required, `debug.PrintStack()` was used to print the stack trace as in the previous case.  
> Generally, in middleware, logs that include a stack trace are commonly kept for high-severity situations like panic so that developers are quickly notified.

Using the `debug.Stack()` method allows for direct handling of the stack trace without outputting it to Stderr.

## Use Cases
`log.Fatal()` internally calls `os.Exit(1)`.

It was designed to immediately terminate the program with an error code, and thus cannot be recovered with `recover()`.

On the other hand, `panic()` can be recovered using `recover()`.

While it shouldn't happen, if a situation arises where recovery is possible, it's better to use `panic()`. 

In general, it’s advisable to use `panic()` in library functions or functions of specific packages. (If the server crashes due to a library and recovery isn’t possible, the outcome could be disastrous.)

It is better to use `log.Fatal()` in cases like the `main()` function when handling errors at the end.

For example, if an error occurs during the process of loading dependencies that prevents the program from being executed, the module that initializes those dependencies would return an `error`, and `log.Fatal()` would be called from the `main()` function.

Here's a simple structure for reference:
> This is a simplified version. Just for reference.
```go
package main

type Dependencies struct {
	DB *sql.DB
	redis *redis.Client
	...
}

func NewDependencies() (*Dependencies, error) {
    db, err := sql.Open("mysql", "user:password@/dbname")
    if err != nil {
        return nil, err
    }

    redis := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    if err = redis.Ping().Err(); err != nil {
        return nil, err
    }

    return &Dependencies{
        DB: db,
        redis: redis,
    }, nil
}

func main() {
    deps, err := NewDependencies()
    if err != nil {
        log.Fatal(err)
    }
    
	// ...
}
```

## + log.Panic()
The `log` package also includes a function called `log.Panic()`.

This adds logging functionality to `panic()`; when the following code is executed:

```go
package main

import (
	"log"
	"log/slog"
	"runtime/debug"
)

func RunWithPanic() {
	log.Panic("This is a panic error")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Info("Recovered from", "error", r)
			debug.PrintStack()
		}
	}()

	RunWithPanic()
}
```

It outputs:
```shell
2025/02/11 20:23:17 This is a panic error
2025/02/11 20:23:17 INFO Recovered from error="This is a panic error"
goroutine 1 [running]:
runtime/debug.Stack()
	/opt/homebrew/opt/go/libexec/src/runtime/debug/stack.go:26 +0x64
runtime/debug.PrintStack()
	/opt/homebrew/opt/go/libexec/src/runtime/debug/stack.go:18 +0x1c
main.main.func1()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:17 +0x8c
panic({0x100ad24e0?, 0x140000100a0?})
	/opt/homebrew/opt/go/libexec/src/runtime/panic.go:785 +0x124
log.Panic({0x1400010af20?, 0x0?, 0x68?})
	/opt/homebrew/opt/go/libexec/src/log/log.go:432 +0x60
main.RunWithPanic(...)
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:10
main.main()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:21 +0x60
```

Compared to `panic`, you can see it includes logging functionality. It functions the same way as `panic()` but commonly features logging capabilities.