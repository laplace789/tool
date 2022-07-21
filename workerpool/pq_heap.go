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
	pq[i].index = i
	pq[j].index = j
}

func (pq *pqHeap) Push(x interface{}) {
	n := len(*pq)
	item := x.(*TaskParam)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *pqHeap) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
