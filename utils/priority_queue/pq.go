package priority_queue

import (
	"container/heap"
	"errors"
	"math"
)

// Package pq implements a priority queue data structure on top of container/heap.
// As an addition to regular operations, it allows an update of an items priority,
// allowing the queue to be used in graph search algorithms like Dijkstra's algorithm.
// Computational complexities of operations are mainly determined by container/heap.
// In addition, a map of items is maintained, allowing O(1) lookup needed for priority updates,
// which themselves are O(log n).

// 权重常量
const (
	price_min = 0.0
	price_max = math.MaxFloat64

	time_stamp_min = 0
	time_stamp_max = math.MaxInt64
)

// PriorityQueue represents the queue
type PriorityQueue struct {
	itemHeap *itemHeap
	lookup   map[interface{}]*item
}

// New initializes an empty priority queue.
func NewPriorityQueue() PriorityQueue {
	return PriorityQueue{
		itemHeap: &itemHeap{},
		lookup:   make(map[interface{}]*item),
	}
}

// Len returns the number of elements in the queue.
func (p *PriorityQueue) Len() int {
	return p.itemHeap.Len()
}

// Insert inserts a new element into the queue. No action is performed on duplicate elements.
func (p *PriorityQueue) Insert(v interface{}, time_stamp int64, price float64) {
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
func (p *PriorityQueue) Remove(value interface{}) (interface{}, error) {
	p.UpdatePriority(value, time_stamp_min, price_max)
	return p.Pop()
}

// Pop removes the element with the highest priority from the queue and returns it.
// In case of an empty queue, an error is returned.
func (p *PriorityQueue) Pop() (interface{}, error) {
	if len(*p.itemHeap) == 0 {
		return nil, errors.New("empty queue")
	}

	item := heap.Pop(p.itemHeap).(*item)
	delete(p.lookup, item.value)
	return item.value, nil
}

// UpdatePriority changes the priority of a given item.
// If the specified item is not present in the queue, no action is performed.
func (p *PriorityQueue) UpdatePriority(x interface{}, time_stamp int64, price float64) {
	item, ok := p.lookup[x]
	if !ok {
		return
	}

	item.time_stamp = time_stamp
	item.price = price
	heap.Fix(p.itemHeap, item.index)
}

type itemHeap []*item

type item struct {
	value    interface{}
	time_stamp int64    // The priority of the item in the queue.
	price float64	 // same as up
	index    int
}

func (ih *itemHeap) Len() int {
	return len(*ih)
}

func (ih *itemHeap) Less(i, j int) bool {
	//return (*ih)[i].priority < (*ih)[j].priority
	// 价格高优先
	if ((*ih)[i].price > (*ih)[j].price) {
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

func (ih *itemHeap) Swap(i, j int) {
	(*ih)[i], (*ih)[j] = (*ih)[j], (*ih)[i]
	(*ih)[i].index = i
	(*ih)[j].index = j
}

func (ih *itemHeap) Push(x interface{}) {
	it := x.(*item)
	it.index = len(*ih)
	*ih = append(*ih, it)
}

func (ih *itemHeap) Pop() interface{} {
	old := *ih
	item := old[len(old)-1]
	*ih = old[0 : len(old)-1]
	return item
}

