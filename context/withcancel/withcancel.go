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
