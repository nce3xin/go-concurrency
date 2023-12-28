# Context上下文

上下文就是指，在 API 之间或者方法调用之间，所传递的除了业务参数之外的额外信息。

Go 标准库中的 Context 常用来提供超时（Timeout）和取消（Cancel）的机制。

## When to use context (Some use cases)

- To pass data to the downstream. Eg. a HTTP request creates a request_id, request_user which needs to be passed around to all downstream functions for distributed tracing.
- When you want to halt the operation in the midway – A HTTP request should be stopped because the client disconnected
- When you want to halt the operation within a specified time from start i.e with timeout – Eg- a HTTP request should be completed in 2 sec or else should be aborted.
- When you want to halt an operation before a certain time – Eg. A cron is running that needs to be aborted in 5 mins if not completed.


## Context 接口

```
type Context interface {
    // It retures a channel when a context is cancelled, timesout (either when deadline is reached or timeout time has finished)
    Done() <-chan struct{}

    // Err will tell why this context was cancelled. A context is cancelled in three scenarios.
    // 1. With explicit cancellation signal
    // 2. Timeout is reached
    // 3. Deadline is reached
    Err() error

    // Used for handling deallines and timeouts
    Deadline() (deadline time.Time, ok bool)

    // Used for passing request scope values
    Value(key interface{}) interface{}
}
```

## 创建新的context

- context.Background()：**可以无脑用这个**。其实这个和context.TODO()的实现是一模一样的。
- context.TODO()

## Context Tree

Before understanding Context Tree please make sure that it is implicitly created in the background when using context. You will find no mention of in go context package itself.

### 1. Two level tree

```
rootCtx := context.Background()
childCtx := context.WithValue(rootCtx, "msgId", "someMsgId")
```

- rootCtx is the empty Context with no functionality
- childCtx is derived from rootCtx and has the functionality of storing request-scoped values. In above example it is storing key-value pair of  {"msgId" : "someMsgId"}

### 2. Three level tree

```
rootCtx := context.Background()
childCtx := context.WithValue(rootCtx, "msgId", "someMsgId")
childOfChildCtx, cancelFunc := context.WithCancel(childCtx)
```

- rootCtx is the empty Context with no functionality
- childCtx is derived from rootCtx and has the functionality of storing request-scoped values. In above example it is storing key-value pair of  {"msgId" : "someMsgId"}
- childOfChildCtx is derived from childCtx . It has the functionality of storing request-scoped values and also it has the functionality of triggering cancellation signals. cancelFunc can be used to trigger cancellation signals

## Deriving From Context

A derived context is can be created in 4 ways

- Passing request-scoped values  -  using WithValue() function of context package
- With cancellation signals - using WithCancel() function of context package
- With deadlines - using WithDeadine() function of context package
- With timeouts - using WithTimeout() function of context package

### 1. context.WithValue()

Used for passing request-scoped values. 

```
withValue(parent Context, key, val interface{}) (ctx Context)
```

用法：

```
#Root Context
ctxRoot := context.Background() - #Root context 

#Below ctxChild has acess to only one pair {"a":"x"}
ctxChild := context.WithValue(ctxRoot, "a", "x") 

#Below ctxChildofChild has access to both pairs {"a":"x", "b":"y"} as it is derived from ctxChild
ctxChildofChild := context.WithValue(ctxChild, "b", "y") 
```

例子：

```
package main

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

func main() {
	helloWorldHandler := http.HandlerFunc(HelloWorld)
	http.Handle("/welcome", injectMsgID(helloWorldHandler))
	_ = http.ListenAndServe(":80", nil)
}

// HelloWorld hello world handler
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	msgID := ""
	if m := r.Context().Value("msgId"); m != nil {
		if value, ok := m.(string); ok {
			msgID = value
		}
	}
	w.Header().Add("msgId", msgID)
	_, _ = w.Write([]byte("Hello, world"))
}

func injectMsgID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msgID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "msgId", msgID)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
```

curl测试：

```shell
# If failed, check the firewall rules, open 80 port
curl -v http://localhost/welcome
```

HTTP response:

```shell
nce3x@zx MINGW64 /d/code/src/go-concurrency/context/withvalue (main)
$ curl -v http://localhost/welcome
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying [::1]:80...
* Connected to localhost (::1) port 80
> GET /welcome HTTP/1.1
> Host: localhost
> User-Agent: curl/8.4.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Msgid: b124e50f-fe1f-4045-a116-083a6d11e22a  <============ Msgid
< Date: Thu, 28 Dec 2023 15:46:59 GMT
< Content-Length: 12
< Content-Type: text/plain; charset=utf-8
<
{ [12 bytes data]
100    12  100    12    0     0  15018      0 --:--:-- --:--:-- --:--:-- 12000Hello, world
* Connection #0 to host localhost left intact
```

### 2. context.WithCancel()

Used for cancellation signals. Below is the signature of WithCancel() function

```
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
```

例子：

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	go task(cancelCtx)
	time.Sleep(time.Second * 3)
	cancelFunc()
	time.Sleep(time.Second * 1)
}

func task(ctx context.Context) {
	i := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Gracefully exit\n")
			fmt.Printf("Ctx error: %v\n", ctx.Err())
			return
		default:
			fmt.Printf("%d\n", i)
			time.Sleep(time.Second * 1)
			i++
		}
	}
}
```

输出：

```
1
2
3
Gracefully exit
context canceled
```

task function will gracefully exit once the cancelFunc is called. Once the cancelFunc is called, the error string is set to "context cancelled" by the context package. That is why the output of ctx.Err() is "context cancelled".

### 3. context.WithTimeout()

Used for time-based cancellation. The signature of the function is

```
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

例子：

```
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	go task(cancelCtx)
	time.Sleep(time.Second * 4)
}

func task(ctx context.Context) {
	i := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Gracefully exit")
			fmt.Println(ctx.Err())
			return
		default:
			fmt.Println(i)
			time.Sleep(time.Second * 1)
			i++
		}
	}
}
```

输出：

```
1
2
3
Gracefully exit
context deadline exceeded
```

task function will gracefully exit once the timeout of 3 seconds is completed. The error string is set to "context deadline exceeded" by the context package. That is why the output of ctx.Err() is "context deadline exceeded".

### 4. context.WithDeadline()

Used for deadline-based cancellation. The signature of the function is

```
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
```

例子：

```
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	cancelCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*5))
	defer cancel()
	go task(cancelCtx)
	time.Sleep(time.Second * 6)
}

func task(ctx context.Context) {
	i := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Gracefully exit")
			fmt.Println(ctx.Err())
			return
		default:
			fmt.Println(i)
			time.Sleep(time.Second * 1)
			i++
		}
	}
}
```

输出：

```
1
2
3
4
5
Gracefully exit
context deadline exceeded
```

task function will gracefully exit once the timeout of 5 seconds is completed as we gave the deadline of Time.now() + 5 seconds. The error string is set to "context deadline exceeded" by the context package. That is why the output of ctx.Err() is "context deadline exceeded".


## 必看文章

- [https://golangbyexample.com/using-context-in-golang-complete-guide/](https://golangbyexample.com/using-context-in-golang-complete-guide/)