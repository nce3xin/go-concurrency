package main

import (
	"fmt"
	"sync"
)

func main() {
	count := 0
	var mu sync.Mutex
	// 用waitGroup等待10个goroutine完成
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				mu.Lock()
				count++
				mu.Unlock()
			}
		}()
	}
	// 等待10个goroutine完成
	wg.Wait()
	fmt.Printf("Count: %d\n", count)
}
