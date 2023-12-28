package main

import "sync"

type Queue struct {
	sync.Mutex
	data []interface{}
}

func NewQueue(n int) *Queue {
	return &Queue{data: make([]interface{}, n)}
}

func (q *Queue) Push(v interface{}) {
	q.Lock()
	defer q.Unlock()
	q.data = append(q.data, v)
}

func (q *Queue) Pop() interface{} {
	q.Lock()
	defer q.Unlock()
	if len(q.data) == 0 {
		return nil
	}
	v := q.data[0]
	q.data = q.data[1:]
	return v
}
