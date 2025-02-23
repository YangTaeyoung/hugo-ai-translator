---
title: Understanding the Difference between log.Fatal() and panic() in Golang
type: blog
date: 2025-02-11
comments: true
---

## "In this case, using `panic()` seems to be better than `log.Fatal()`"
I recently received feedback while using `log.Fatal()`.

Hmm? Isn't `log.Fatal()` just better at logging? I thought to myself.

Embarrassingly, it wasn't until recently that I clearly understood the difference between `log.Fatal()` and `panic()` in Golang. So, I'm going to try to organize that now.

## Difference between log.Fatal() and panic()
Both `log.Fatal()` and `panic()` are functions that terminate the program. Let's see how they work in code.

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

Running the above code will produce the following result:

```shell
2025/02/11 20:02:31 This is a fatal error
```

Now let's look at the code using `panic()`.

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

Running the above code will produce the following result:

```shell
panic: This is a panic error

goroutine 1 [running]:
main.RunWithPanic(...)
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:8
main.main()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:12 +0x30
```

Looking at the code, `log.Fatal()` prints the error and terminates the program, while `panic()` also prints the error and terminates the program but additionally prints a stack trace.

### Recovering using recover()
When using `panic()`, the program terminates, but using `recover()` allows for recovery without program termination.

Ideally, there should be no panic, but developers are human and mistakes happen. Therefore, in places like API servers, middleware is often created to `recover()` from `panic()` situations, preventing the server from unexpectedly crashing.

To understand the difference clearly, let's first try to recover from `log.Fatal()`.

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

Running the above code will produce the following result:

```shell
2025/02/11 20:07:49 This is a fatal error
```

There is no recovery. Now let's try to recover from `panic()`.

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

Running the above code will produce the following result:

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

Although not necessary, `debug.PrintStack()` is used to output stack traces similar to before.

## Purpose
`log.Fatal()` internally calls `os.Exit(1)`.

It was designed to immediately terminate the program with an error code, hence it cannot be recovered with `recover()`.

On the other hand, `panic()` can be recovered with `recover()`.

If a situation arises where recovery is possible in case of errors, it is better to use `panic()`. Typically, it is advisable to use `panic()` for library functions or specific package functions. (It would be disastrous if the server crashed due to a library and could not be recovered.)

In the case of `log.Fatal()`, it is best used when handling errors at the end, such as in the `main()` function.

For instance, when an error occurs that prevents the program from running during initialization of dependencies, the module that initializes the dependencies returns an `error`, and the `main()` function uses `log.Fatal()`.

A sample structure might look like this:
> This is a simplified version. Read it for reference only.
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

## + `log.Panic()`
There is also the function `log.Panic()` in the `log` package.

By running the following code, which adds logging functionality to `panic()`:

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

The following output will be generated:
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

Compared to panic, it is a version that adds logging functionality. Although it raises `panic()` in the same way, it logs the error.

## Reference
- https://pkg.go.dev/log#Fatal
- Courtesy of feedback from the team leader during code reviews
