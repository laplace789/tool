package workerpool

import (
	"container/heap"
	"sync"
)

type PriorityQueue interface {
	EnQueue(*TaskParam)  //add Task to priority queue
	DeQueue() *TaskParam //return Max Priority task
	Len() int
}

// PriorityQueueImp Todo thread safe
type PriorityQueueImp struct {
	mux sync.RWMutex
	pq  *pqHeap
}

func NewPriorityQueue() PriorityQueue {
	hp := make(pqHeap, 0)
	heap.Init(&hp)

	pq := new(PriorityQueueImp)
	pq.pq = &hp
	return pq
}

func (p *PriorityQueueImp) EnQueue(param *TaskParam) {
	p.mux.Lock()
	defer p.mux.Unlock()
	heap.Push(p.pq, param)
}

func (p *PriorityQueueImp) DeQueue() *TaskParam {
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.pq.Len() == 0 {
		return nil
	}
	return heap.Pop(p.pq).(*TaskParam)
}

func (p *PriorityQueueImp) Len() int {
	p.mux.RLock()
	defer p.mux.RUnlock()
	return p.pq.Len()
}
