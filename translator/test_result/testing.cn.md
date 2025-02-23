---
title: 了解Golang中log.Fatal()与panic()的区别
type: blog
date: 2025-02-11
comments: true
---

## "在这种情况下，我觉得使用`panic()`比`log.Fatal()`更好"
最近在使用`log.Fatal()`时，听到了上述反馈。

嗯？`log.Fatal()`不是更好地记录`log`吗？我想道

尴尬的事实是，直到最近我才清楚了Golang中`log.Fatal()`和`panic()`的区别，所以我想在这个机会来整理一下。

## log.Fatal()和panic()的区别
`log.Fatal()`和`panic()`都是用于终止程序的函数。让我们通过代码来看看它们之间的区别。

```go
package main

import (
    "log"
	"log/slog"
)

func RunWithFatal() {
	log.Fatal("这是一个致命错误")
}

func main() {
	RunWithFatal()

	slog.Info("这不会执行")
}
```

当运行上面的代码时，会看到以下结果。

```shell
2025/02/11 20:02:31 这是一个致命错误
```

现在让我们看看使用`panic()`的代码。

```go
package main

import (
    "log/slog"
)

func RunWithPanic() {
    panic("这是一个panic错误")
}

func main() {
    RunWithPanic()

    slog.Info("这不会执行")
}
```

当运行上面的代码时，会看到以下结果。

```shell
panic: 这是一个panic错误

goroutine 1 [running]:
main.RunWithPanic(...)
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:8
main.main()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:12 +0x30
```

从上面的代码可以看出，`log.Fatal()`会输出错误并终止程序，而`panic()`也会输出错误并终止程序，但`panic()`还会输出堆栈跟踪。

### 使用recover()进行恢复
使用`panic()`会导致程序终止，但使用`recover()`可以使程序不终止而是恢复。

事实上，最好根本不要出现panic，但作为开发人员，我们不可避免地会犯错误。
因此，通常在API服务器等地方，会创建一个中间件来恢复`panic()`，以防止服务器意外中断的情况发生。

为了更清楚了解区别，让我们先来看看如何在log.Fatal()中使用recover()
```go
package main

import (
    "log"
    "log/slog"
)

func RunWithFatal() {
    log.Fatal("这是一个致命错误")
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            slog.Info("从错误中恢复", "error", r)
        }
    }()

    RunWithFatal()

    slog.Info("这不会执行")
}
```

当运行上面的代码时，会看到以下结果。

```shell
2025/02/11 20:07:49 这是一个致命错误
```

没有恢复。现在让我们来看看如何在`panic()`中使用`recover()`

```go
package main

import (
	"log/slog"
	"runtime/debug"
)

func RunWithPanic() {
	panic("这是一个panic错误")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Info("从错误中恢复", "error", r)
			debug.PrintStack()
		}
	}()

	RunWithPanic()
}
```

当运行上面的代码时，会看到以下结果。

```shell
2025/02/11 20:09:51 INFO 从错误中恢复 error="这是一个panic错误"
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

虽然不使用recover()也可以，但为了输出与之前相同的堆栈跟踪信息，我们使用了`debug.PrintStack()`。

## 用途
`log.Fatal()`内部调用了`os.Exit(1)`。

它最初是设计用于立即结束程序，因此无法通过`recover()`来恢复。

另一方面，`panic()`可以通过`recover()`来恢复。

如果出现不应该出现的情况，但在出现时可以进行恢复，那么最好使用`panic()`。

通常在库函数或特定包函数的情况下，最好使用`panic()`。（如果由于库而导致服务器崩溃并且无法恢复，结果将是悲惨的。）

在`main()`函数等地方，当最终要处理`error`时，最好使用`log.Fatal()`。

举例来说，在初始化依赖项时由于无法启动程序而出现错误时，初始化这些依赖项的模块会返回`error`，而`main()`函数会调用`log.Fatal()`。

如果看其结构会是这样的。
> 以下只是一个简单示例，仅供参考。
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
`log`包还有一个名为`log.Panic()`的函数。

它是`panic()`加上日志功能，执行以下代码将会输出

```go
package main

import (
	"log"
	"log/slog"
	"runtime/debug"
)

func RunWithPanic() {
	log.Panic("这是一个panic错误")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Info("从错误中恢复", "error", r)
			debug.PrintStack()
		}
	}()

	RunWithPanic()
}

```

将输出如下。
```shell
2025/02/11 20:23:17 这是一个panic错误
2025/02/11 20:23:17 INFO 从错误中恢复 error="这是一个panic错误"
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

与panic相比，此函数加入了日志功能。它可以被看作是触发`panic()`时同时记录日志的特性。

## 参考资料
- https://pkg.go.dev/log#Fatal
- 由总监进行的代码审查