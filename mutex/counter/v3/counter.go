package main

import (
	"fmt"
	"sync"
)

// Counter 线程安全的计数器类型
type Counter struct {
	Name string

	// 如果嵌入的 struct 有多个字段，我们一般会把 Mutex 放在要控制的字段上面，然后使用空格把字段分隔开来
	sync.Mutex
	Counter int
}

// Incr 加1的方法，内部使用互斥锁保护
func (c *Counter) Incr() {
	c.Lock()
	c.Counter++
	c.Unlock()
}

func main() {
	var counter Counter
	var wg sync.WaitGroup
	wg.Add(10)
	// 启动10个goroutine
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			// 执行10万次累加
			for j := 0; j < 100000; j++ {
				// 受到锁保护的方法
				counter.Incr()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Count: %d\n", counter.Counter)
}
