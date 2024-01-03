package main

import (
	"fmt"
	"sync"
)

type Foo struct {
	first  chan struct{}
	second chan struct{}
	third  chan struct{}
}

func NewFoo() *Foo {
	return &Foo{
		first:  make(chan struct{}),
		second: make(chan struct{}),
		third:  make(chan struct{}),
	}
}

func (f *Foo) First() {
	<-f.first
	fmt.Printf("first ")
	f.second <- struct{}{}
}

func (f *Foo) Second() {
	<-f.second
	fmt.Printf("second ")
	f.third <- struct{}{}
}

func (f *Foo) Third() {
	<-f.third
	fmt.Printf("third ")
}

// expected result: first second third
func main() {
	f := NewFoo()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		f.First()
	}()

	go func() {
		defer wg.Done()
		f.Second()
	}()

	go func() {
		defer wg.Done()
		f.Third()
	}()

	f.first <- struct{}{}

	wg.Wait()
}
