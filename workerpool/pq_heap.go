package workerpool

// A pqHeap implements heap.Interface and holds Items.
type pqHeap []*TaskParam

func (pq pqHeap) Len() int { return len(pq) }

func (pq pqHeap) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].TaskPriority > pq[j].TaskPriority
}

func (pq pqHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *pqHeap) Push(x interface{}) {
	item := x.(*TaskParam)
	*pq = append(*pq, item)
}

func (pq *pqHeap) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
