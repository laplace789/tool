package slice

import "sort"

// NewPriority return a slice of PriorityItem
func NewPriority[V any](lengthAndCap ...int) *Priority[V] {
	p := &Priority[V]{}
	if len(lengthAndCap) > 0 {
		var length = lengthAndCap[0]
		var c int
		if len(lengthAndCap) > 1 {
			c = lengthAndCap[1]
		}
		p.items = make([]*PriorityItem[V], length, c)
	}
	return p
}

// Priority  is a slice order by priority
type Priority[V any] struct {
	items []*PriorityItem[V]
}

func (p *Priority[V]) Len() int {
	return len(p.items)
}

func (p *Priority[V]) Cap() int {
	return cap(p.items)
}

// Clear  slice
func (p *Priority[V]) Clear() {
	p.items = p.items[:0]
}

// Append add single priorityItem to slice and sort it
func (p *Priority[V]) Append(value V, priority int) {
	p.items = append(p.items, NewPriorityItem(value, priority))
	p.sort()
}

// Appends add multi priorityItem into slice
func (p *Priority[V]) Appends(priority int, vs ...V) {
	for _, v := range vs {
		p.Append(v, priority)
	}
	p.sort()
}

// Get return PriorityItem of index
func (p *Priority[V]) Get(index int) *PriorityItem[V] {
	if index-1 > len(p.items) {
		return nil
	}
	return p.items[index]
}

// GetValue get value of items[index]
func (p *Priority[V]) GetValue(index int) V {
	if index-1 > len(p.items) {
		var v V
		return v
	}
	return p.items[index].Value()
}

// GetPriority get priority of items[index]
func (p *Priority[V]) GetPriority(index int) int {
	if index-1 > len(p.items) {
		return -1
	}
	return p.items[index].Priority()
}

func (p *Priority[V]) Set(v V, index, priority int) {
	before := p.items[index]
	p.items[index] = NewPriorityItem(v, priority)
	if before.Priority() != priority {
		p.sort()
	}
}

func (p *Priority[V]) SetValue(index int, value V) {
	p.items[index].value = value
}

func (p *Priority[V]) SetPriority(index, priority int) {
	p.items[index].priority = priority
	p.sort()
}

// Range iterator slice if action return false then stop
func (p *Priority[V]) Range(action func(index int, item *PriorityItem[V]) bool) {
	for i, item := range p.items {
		if !action(i, item) {
			break
		}
	}
}

// RangeValue iterator slice if action return false then stop
func (p *Priority[V]) RangeValue(action func(index int, value V) bool) {
	p.Range(func(index int, item *PriorityItem[V]) bool {
		return action(index, item.Value())
	})
}

// RangePriority iterator slice if action return false then stop
func (p *Priority[V]) RangePriority(action func(index int, priority int) bool) {
	p.Range(func(index int, item *PriorityItem[V]) bool {
		return action(index, item.Priority())
	})
}

// sort  sorting slice using priority
func (p *Priority[V]) sort() {
	if len(p.items) <= 1 {
		return
	}
	//sort the slice
	sort.Slice(p.items, func(i, j int) bool {
		return p.items[i].Priority() < p.items[j].Priority()
	})

	//link prev and next after sort
	for i := 0; i < len(p.items); i++ {
		if i == 0 {
			//first item
			p.items[i].prev = nil
			p.items[i].next = p.items[i+1]
		} else if i == len(p.items)-1 {
			//last item
			p.items[i].prev = p.items[i-1]
			p.items[i].next = nil
		} else {
			p.items[i].prev = p.items[i-1]
			p.items[i].next = p.items[i+1]
		}
	}
}
