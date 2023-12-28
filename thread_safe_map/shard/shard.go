package main

import "sync"

const ShardCount = 32

// ConcurrentMap 分成SHARD_COUNT个分片的map
type ConcurrentMap []*ConcurrentMapShard

// ConcurrentMapShard 通过RWMutex保护的线程安全的分片，包含一个map
type ConcurrentMapShard struct {
	sync.RWMutex
	items map[string]interface{}
}

// NewConcurrentMap 创建并发map
func NewConcurrentMap() ConcurrentMap {
	m := make(ConcurrentMap, ShardCount)
	for i := 0; i < ShardCount; i++ {
		m[i] = &ConcurrentMapShard{items: make(map[string]interface{})}
	}
	return m
}

// GetShard 根据key计算分片索引
func (m *ConcurrentMap) GetShard(key string) *ConcurrentMapShard {
	return (*m)[hash(key)%ShardCount]
}

func hash(k string) uint {
	// TODO
	return 0
}

func (m *ConcurrentMap) Set(k string, v interface{}) {
	// 根据key计算出对应的分片
	shard := m.GetShard(k)
	// 对这个分片加锁，执行业务操作
	shard.Lock()
	defer shard.Unlock()
	shard.items[k] = v
}

func (m *ConcurrentMap) Get(k string) (interface{}, bool) {
	// 根据key计算出对应的分片
	shard := (*m).GetShard(k)
	shard.RLock()
	defer shard.Unlock()
	// 从这个分片读取key的值
	v, ok := shard.items[k]
	return v, ok
}
