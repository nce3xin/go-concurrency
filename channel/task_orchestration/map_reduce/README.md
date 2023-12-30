# 单机单进程的 map reduce

map-reduce 是一种处理数据的方式，最早是由 Google 公司研究提出的一种面向大规模数据处理的并行计算模型和方法，开源的版本是 hadoop，前几年比较火。

不过，这里说的并不是分布式的 map-reduce，而是单机单进程的 map-reduce 方法。

map-reduce 分为两个步骤，第一步是映射（map），处理队列中的数据，第二步是规约（reduce），把列表中的每一个元素按照一定的处理方式处理成结果，放入到结果队列中。

就像做汉堡一样，map 就是单独处理每一种食材，reduce 就是从每一份食材中取一部分，做成一个汉堡。

## map的逻辑

```
func mapChan(in <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
    out := make(chan interface{}) //创建一个输出chan
    if in == nil { // 异常检查
        close(out)
        return out
    }

    go func() { // 启动一个goroutine,实现map的主要逻辑
        defer close(out)
        for v := range in { // 从输入chan读取数据，执行业务操作，也就是map操作
            out <- fn(v)
        }
    }()

    return out
}
```

## reduce的逻辑

```
func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
    if in == nil { // 异常检查
        return nil
    }

    out := <-in // 先读取第一个元素
    for v := range in { // 实现reduce的主要逻辑
        out = fn(out, v)
    }

    return out
}
```