# WaitGroup 协同等待，任务编排利器

它要解决的就是并发 - 等待的问题：现在有一个 goroutine A 在检查点（checkpoint）等待一组 goroutine 全部完成，如果在执行任务的这些 goroutine 还没全部完成，那么 goroutine A 就会阻塞在检查点，直到所有 goroutine 都完成后才能继续执行。

## 基本用法

Go 标准库中的 WaitGroup 提供了三个方法：

```
func (wg *WaitGroup) Add(delta int)
func (wg *WaitGroup) Done()
func (wg *WaitGroup) Wait()
```

- Add，用来设置 WaitGroup 的计数值；
- Done，用来将 WaitGroup 的计数值减 1，其实就是调用了 Add(-1)；
- Wait，调用这个方法的 goroutine 会一直阻塞，直到 WaitGroup 的计数值变为 0。