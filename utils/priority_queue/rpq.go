package priority_queue

import (
	"container/heap"
	"errors"
)

// Package pq implements a priority queue data structure on top of container/heap.
// As an addition to regular operations, it allows an update of an items priority,
// allowing the queue to be used in graph search algorithms like Dijkstra's algorithm.
// Computational complexities of operations are mainly determined by container/heap.
// In addition, a map of items is maintained, allowing O(1) lookup needed for priority updates,
// which themselves are O(log n).

// RPriorityQueue represents the queue
type RPriorityQueue struct {
	itemHeap *itemHeapR
	lookup   map[interface{}]*item
}

// New initializes an empty priority queue.
func NewRPriorityQueue() RPriorityQueue {
	return RPriorityQueue{
		itemHeap: &itemHeapR{},
		lookup:   make(map[interface{}]*item),
	}
}

// Len returns the number of elements in the queue.
func (p *RPriorityQueue) Len() int {
	return p.itemHeap.Len()
}

// Insert inserts a new element into the queue. No action is performed on duplicate elements.
func (p *RPriorityQueue) Insert(v interface{}, time_stamp int64, price float64) {
	_, ok := p.lookup[v]
	if ok {
		return
	}

	newItem := &item{
		value:    v,
		time_stamp: time_stamp,
		price: price,
	}
	heap.Push(p.itemHeap, newItem)
	p.lookup[v] = newItem
}

// remove 删除一个元素
func (p *RPriorityQueue) Remove(value interface{}) (interface{}, error) {
	p.UpdatePriority(value, time_stamp_min, price_min)
	return p.Pop()
}

// Pop removes the element with the highest priority from the queue and returns it.
// In case of an empty queue, an error is returned.
func (p *RPriorityQueue) Pop() (interface{}, error) {
	if len(*p.itemHeap) == 0 {
		return nil, errors.New("empty queue")
	}

	item := heap.Pop(p.itemHeap).(*item)
	delete(p.lookup, item.value)
	return item.value, nil
}

// UpdatePriority changes the priority of a given item.
// If the specified item is not present in the queue, no action is performed.
func (p *RPriorityQueue) UpdatePriority(x interface{}, time_stamp int64, price float64) {
	item, ok := p.lookup[x]
	if !ok {
		return
	}

	item.time_stamp = time_stamp
	item.price = price
	heap.Fix(p.itemHeap, item.index)
}

type itemHeapR []*item

func (ih *itemHeapR) Len() int {
	return len(*ih)
}

func (ih *itemHeapR) Less(i, j int) bool {
	//return (*ih)[i].priority < (*ih)[j].priority
	// 价格低优先
	if ((*ih)[i].price < (*ih)[j].price) {
		return true
	}

	// 时间优先
	if ((*ih)[i].price == (*ih)[j].price) {
		if ((*ih)[i].time_stamp < (*ih)[j].time_stamp) {
			return true
		}
	}

	return false
}

func (ih *itemHeapR) Swap(i, j int) {
	(*ih)[i], (*ih)[j] = (*ih)[j], (*ih)[i]
	(*ih)[i].index = i
	(*ih)[j].index = j
}

func (ih *itemHeapR) Push(x interface{}) {
	it := x.(*item)
	it.index = len(*ih)
	*ih = append(*ih, it)
}

func (ih *itemHeapR) Pop() interface{} {
	old := *ih
	item := old[len(old)-1]
	*ih = old[0 : len(old)-1]
	return item
}

