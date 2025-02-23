---
title: Golang에서 log.Fatal() vs panic() 차이점에 대해 알아보기 
type: blog
date: 2025-02-11
comments: true
---

## "이런 경우 `log.Fatal()`보다 `panic()`을 사용하는게 더 좋은 것 같아요"
최근 `log.Fatal()`을 사용하다 위와 같은 피드백을 들었다.

음? `log.Fatal()`은 `log`를 좀 더 잘 찍어주는거 아닌가? 하고 생각했다

부끄러운 사실이지만, 최근이 되어서야 Golang에서 `log.Fatal()`과 `panic()`의 차이점을 명확하게 알게 되어 이번 기회에 정리해보려고 한다.

## log.Fatal()과 panic()의 차이점
`log.Fatal()`과 `panic()`은 둘 다 프로그램을 종료시키는 함수이다. 코드로 동작을 살펴보자

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

위 코드를 실행하면 다음과 같은 결과를 볼 수 있다.

```shell
2025/02/11 20:02:31 This is a fatal error
```

이제 `panic()`을 사용한 코드를 살펴보자.

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

위 코드를 실행하면 다음과 같은 결과를 볼 수 있다.

```shell
panic: This is a panic error

goroutine 1 [running]:
main.RunWithPanic(...)
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:8
main.main()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:12 +0x30
```

위 코드를 보면 `log.Fatal()`은 에러를 출력하고 프로그램을 종료시키지만, `panic()`은 에러를 출력하고 프로그램을 종료시키는 것은 동일하지만, `panic()`은 스택 트레이스를 출력한다.

### recover()를 사용하여 복구하기
`panic()`을 사용하면 프로그램이 종료되지만, `recover()`를 사용하면 프로그램을 종료하지 않고 복구할 수 있다.

사실상 panic이 아예 없다면 좋겠지만, 개발자도 사람이기에 실수를 할 수 밖에 없다.
그렇기에 일반적으로 API 서버와 같은 곳에서는 `panic()`을 `recover()`하는 미들웨어를 만들어 복구하는 방식으로 서버가 의도치 않게 죽어버리는 상황을 방지하곤 한다.

차이를 명확하게 알기 위해 먼저 log.Fatal()을 recover()를 해보자
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

위 코드를 실행하면 다음과 같은 결과를 볼 수 있다.

```shell
2025/02/11 20:07:49 This is a fatal error
```

복구되지 않은 것이다. 그렇다면 `panic()`을 `recover()`를 해보자

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

위 코드를 실행하면 다음과 같은 결과를 볼 수 있다.

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

사용하지 않아도 되지만, 좀전과 같은 스택 트레이스를 출력하기 위해 `debug.PrintStack()`를 사용하였다.
> 일반적으로 미들웨어에서는 심각도가 높은 panic과 같은 상황에서 개발자가 빨리 알아차릴 수 있도록 위와 같은 스택트레이스를 포함한 로그를 많이 남긴다.

`debug.Stack()` 메서드를 통해 Stderr에 출력하지 않고도, 스택 트레이스를 직접 핸들링할 수도 있다.

## 용도
`log.Fatal()`은 내부적으로 `os.Exit(1)`을 호출한다.

애초에 에러코드와 함께 프로그램을 즉시 종료하기 위해 만들어졌고, 그렇기에 `recover()`로 복구할 수 없다.

반면 `panic()`은 `recover()`로 복구할 수 있다.

발생해선 안되지만, 발생했을 때 복구할 수 있는 상황이라면 `panic()`을 사용하는 것이 더 좋다.

일반적으로 라이브러리 함수나, 특정 패키지의 함수의 경우에 `panic()`을 사용하는 것이 좋다. (라이브러리 때문에 서버가 죽었는데 복구도 되지 않는다면 그 결과는 참혹할 것이다.)

`log.Fatal()`의 경우에는 `main()` 함수와 같은 경우에 `error`를 최종적으로 핸들링할 때 사용하는 것이 좋다.

예를 들어 의존성을 불러오는 과정에서 프로그램을 실행시킬 수 없을만한 에러가 발생했을 때 해당 의존성 초기화 하는 모듈은 `error`를 반환하고, `main()` 함수에서 `log.Fatal()`을 호출하는 식이다.

구조를 본다면 이럴 것이다.
> 단순하게 적은 것이다. 참고용으로만 보자.
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
`log` 패키지에는 `log.Panic()`이라는 함수도 있다.

`panic()`에 로그 기능을 더한 것으로 아래와 같은 코드를 실행시키면

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

다음과 같이 출력된다.
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

panic과 비교했을 때, 로그 기능을 더한 것이라 보면 된다. 동일하게 `panic()`을 발생시키지만, 로그를 남기는 것이 특징이라고 볼 수 있다.

## Reference
- https://pkg.go.dev/log#Fatal
- 코드 리뷰 해주신 대표님