---
title: Understanding the Differences between log.Fatal() and panic() in Golang
type: blog
date: 2025-02-11
comments: true
---

## "In such cases, it seems better to use `panic()` than `log.Fatal()`"
I recently heard this feedback while using `log.Fatal()`.

Huh? Doesn't `log.Fatal()` just log a bit better? I thought.

Embarrassingly, it wasn't until recently that I clearly understood the differences between `log.Fatal()` and `panic()` in Golang. So, I wanted to take this opportunity to organize it.

## Differences between log.Fatal() and panic()
Both `log.Fatal()` and `panic()` are functions that terminate the program. Let's see how they work in code:

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

When you run the above code, you will see the following result:

```shell
2025/02/11 20:02:31 This is a fatal error
```

Now, let's look at the code using `panic()`.

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

When you run the above code, you will see the following result:

```shell
panic: This is a panic error

goroutine 1 [running]:
main.RunWithPanic(...)
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:8
main.main()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:12 +0x30
```

From the code, it can be observed that while `log.Fatal()` outputs an error and terminates the program, `panic()` outputs an error, also terminates the program, and additionally prints a stack trace.

### Recovering using recover()
When using `panic()`, the program terminates, but with `recover()`, you can recover without ending the program.

Ideally, there should be no panics, but developers are prone to mistakes. Therefore, in places like API servers, a middleware is created to `recover()` from `panic()`, preventing the server from unexpectedly crashing.

To understand the difference clearly, let's first recover from `log.Fatal()`.
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

When you run the above code, you will see the following result:

```shell
2025/02/11 20:07:49 This is a fatal error
```

It remains unrecovered. Now, let's recover from `panic()`.

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

When you run the above code, you will see the following result:

```shell
2025/02/11 20:09:51 INFO Recovered from error="This is a panic error"
goroutine 1 [running]:
runtime/debug.Stack()
	/opt/homebrew/opt/go/libexec/src/runtime/debug/stack.go:26 +0x64
runtime/debug.PrintStack()
	/opt/homebrew/opt/go/libexec/src/runtime/debug/stack.go:18 +0x1c
main.main.func1()
	/Users/code_kir...
```

Although not necessary, `debug.PrintStack()` is used to output a stack trace as seen previously.

## Purpose
`log.Fatal()` internally calls `os.Exit(1)`.

It was created to immediately terminate the program with an error code, making it irrecoverable with `recover()`.

On the other hand, `panic()` can be recovered with `recover()`.

If there's a situation where an error should not occur but if it does, it can be recovered from, using `panic()` is better.

Generally, for library functions or specific package functions, it is advisable to use `panic()`. (If a server crashes due to a library and cannot be recovered, it can have catastrophic consequences.)

For `log.Fatal()`, it is recommended to use it when handling errors finally in the `main()` function.

For instance, in the process of loading dependencies if an error occurs that prevents the program from running, the module initializing the dependency returns an `error`, and in the `main()` function, `log.Fatal()` is called.

The structure might look like this:
> This is a simplification. Please use it for reference only.
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
In the `log` package, there is a function called `log.Panic()`.

It is an extension of `panic()` with logging capabilities. When you run code like this:

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

It will output like this:
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
	/opt/homebrew/op...
```

Compared to `panic()`, it includes logging capabilities. It triggers `panic()` but also logs, distinguishing it from a standard `panic()`.

## Reference
- https://pkg.go.dev/log#Fatal
- Code review by the CEO
