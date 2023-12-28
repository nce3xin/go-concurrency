# Once

Once 可以用来执行且仅仅执行一次动作，常常用于单例对象的初始化场景。

sync.Once 只暴露了一个方法 Do，你可以多次调用 Do 方法，但是只有第一次调用 Do 方法时 f 参数才会执行，这里的 f 是一个无参数无返回值的函数。

```
func (o *Once) Do(f func())
```

因为当且仅当第一次调用 Do 方法的时候参数 f 才会执行，即使第二次、第三次、第 n 次调用时 f 参数的值不一样，也不会被执行，比如下面的例子，虽然 f1 和 f2 是不同的函数，但是第二个函数 f2 就不会执行。

```go
package main


import (
    "fmt"
    "sync"
)

func main() {
    var once sync.Once

    // 第一个初始化函数
    f1 := func() {
        fmt.Println("in f1")
    }
    once.Do(f1) // 打印出 in f1

    // 第二个初始化函数
    f2 := func() {
        fmt.Println("in f2")
    }
    once.Do(f2) // 无输出
}
```

因为这里的 f 参数是一个无参数无返回的函数，所以你可能会通过闭包的方式引用外面的参数，比如：

```
var addr = "baidu.com"

var conn net.Conn
var err error

once.Do(func() {
conn, err = net.Dial("tcp", addr)
})
```

而且在实际的使用中，绝大多数情况下，你会使用闭包的方式去初始化外部的一个资源。

重点介绍一下很值得我们学习的 math/big/sqrt.go 中实现的一个数据结构，它通过 Once 封装了一个只初始化一次的值：

```
   // 值是3.0或者0.0的一个数据结构
   var threeOnce struct {
    sync.Once
    v *Float
  }
  
    // 返回此数据结构的值，如果还没有初始化为3.0，则初始化
  func three() *Float {
    threeOnce.Do(func() { // 使用Once初始化
      threeOnce.v = NewFloat(3.0)
    })
    return threeOnce.v
  }
```

它将 sync.Once 和 *Float 封装成一个对象，提供了只初始化一次的值 v。 你看它的 three 方法的实现，虽然每次都调用 threeOnce.Do 方法，但是参数只会被调用一次。

当你使用 Once 的时候，你也可以尝试采用这种结构，将值和 Once 封装成一个新的数据结构，提供只初始化一次的值。

## 总结

总结一下 Once 并发原语解决的问题和使用场景：Once 常常用来初始化单例资源，或者并发访问只需初始化一次的共享资源，或者在测试的时候初始化一次测试资源。
