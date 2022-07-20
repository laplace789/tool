package workerpool

import "container/heap"

type PriorityQueue interface {
	EnQueue(*TaskParam)  //add Task to priority queue
	DeQueue() *TaskParam //return Max Priority task
	Len() int
}

// PriorityQueueImp Todo thread safe
type PriorityQueueImp struct {
	pq *pqHeap
}

func NewPriorityQueue() PriorityQueue {
	hp := make(pqHeap, 0)
	heap.Init(&hp)

	pq := new(PriorityQueueImp)
	pq.pq = &hp
	return pq
}

func (p PriorityQueueImp) EnQueue(param *TaskParam) {
	heap.Push(p.pq, param)
}

func (p PriorityQueueImp) DeQueue() *TaskParam {
	if p.pq.Len() == 0 {
		return nil
	}
	return heap.Pop(p.pq).(*TaskParam)
}

func (p PriorityQueueImp) Len() int {
	return p.pq.Len()
}
