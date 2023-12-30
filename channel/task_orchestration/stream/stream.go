package main

// asStream 创建流的方法。这个方法把一个数据 slice 转换成流
func asStream(done <-chan struct{}, values ...interface{}) <-chan interface{} {
	// 创建一个unbuffered的channel
	c := make(chan interface{})
	// 启动一个goroutine，往s中塞数据
	go func() {
		// 退出时关闭chan
		defer close(c)
		// 遍历数组
		for _, v := range values {
			select {
			case <-done:
				return
			case c <- v: // 将数组元素塞入到chan中
			}
		}
	}()
	return c
}

// takeN 只取流中的前 n 个数据
func takeN(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	// 创建输出流
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		// 只读取前num个元素
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- valueStream: //从输入流中读取元素
			}
		}
	}()
	return takeStream
}
