package main

import (
	"fmt"
	"time"
)

type FooBar struct {
	chs []chan struct{}
}

func NewFooBar() *FooBar {
	return &FooBar{
		chs: []chan struct{}{
			make(chan struct{}),
			make(chan struct{}),
		},
	}
}

func (fb *FooBar) Foo() {
	for {
		<-fb.chs[0]
		fmt.Printf("foo")
		time.Sleep(time.Second * 1)
		fb.chs[1] <- struct{}{}
	}
}

func (fb *FooBar) Bar() {
	for {
		<-fb.chs[1]
		fmt.Printf("bar\n")
		time.Sleep(time.Second * 1)
		fb.chs[0] <- struct{}{}
	}
}

func main() {
	fb := NewFooBar()

	go func() {
		fb.Foo()
	}()

	go func() {
		fb.Bar()
	}()

	fb.chs[0] <- struct{}{}

	select {}
}
