package main

import "fmt"

// 使用 map-reduce 模式处理一组整数，map 函数就是为每个整数乘以 10，reduce 函数就是把 map 处理的结果累加起来

func mapChan(in <-chan interface{}, fn func(v interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{})
	if in == nil {
		close(out)
		return out
	}
	go func() {
		defer close(out)
		for v := range in {
			out <- fn(v)
		}
	}()
	return out
}

func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
	if in == nil {
		return nil
	}
	out := <-in
	for v := range in {
		out = fn(out, v)
	}
	return out
}

// 生成一个数据流
func asStream(done <-chan struct{}) <-chan interface{} {
	c := make(chan interface{})
	values := []int{1, 2, 3, 4, 5}
	go func() {
		defer close(c)
		for _, v := range values {
			select {
			case <-done:
				return
			case c <- v:
			}
		}
	}()
	return c
}

func main() {
	in := asStream(nil)
	mapFn := func(v interface{}) interface{} {
		return v.(int) * 10
	}
	reduceFn := func(r, v interface{}) interface{} {
		return r.(int) + v.(int)
	}
	sum := reduce(mapChan(in, mapFn), reduceFn)
	fmt.Printf("Sum: %d\n", sum)
}
