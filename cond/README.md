# Cond条件变量

用于实现等待-通知（wait-notify）机制。

但Cond 在实际项目中被使用的机会比较少，原因有2个：

- 第一，同样的场景我们会使用其他的并发原语来替代。Go 特有的 Channel 类型，有一个应用很广泛的模式就是通知机制，这个模式使用起来也特别简单。所以很多情况下，我们会使用 Channel 而不是 Cond 实现 wait/notify 机制。
- 第二，对于简单的 wait/notify 场景，比如等待一组 goroutine 完成之后继续执行余下的代码，我们会使用 WaitGroup 来实现。因为 WaitGroup 的使用方法更简单，而且不容易出错。比如，上面百米赛跑的问题，就可以很方便地使用 WaitGroup 来实现。