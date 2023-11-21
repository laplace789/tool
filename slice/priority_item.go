package slice

func NewPriorityItem[V any](v V, priority int) *PriorityItem[V] {
	return &PriorityItem[V]{value: v, priority: priority}
}

type PriorityItem[V any] struct {
	next     *PriorityItem[V]
	prev     *PriorityItem[V]
	value    V
	priority int //priority
}

func (p *PriorityItem[V]) Value() V {
	return p.value
}

func (p *PriorityItem[V]) Priority() int {
	return p.priority
}

func (p *PriorityItem[V]) Next() *PriorityItem[V] {
	return p.next
}

func (p *PriorityItem[V]) Prev() *PriorityItem[V] {
	return p.prev
}
