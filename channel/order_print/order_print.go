package main

import (
	"fmt"
	"time"
)

func main() {
	chs := []chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}

	for i := 0; i < 4; i++ {
		go func(i int) {
			for {
				<-chs[i%4]
				fmt.Printf("%d\n", i+1)
				time.Sleep(time.Second * 1)
				chs[(i+1)%4] <- struct{}{}
			}
		}(i)
	}

	chs[0] <- struct{}{}
	select {}
}
