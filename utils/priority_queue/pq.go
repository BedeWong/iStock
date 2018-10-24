package priority_queue

import "container/heap"

// 根据 官方文档改编的 优先队列：买盘队列
// An Item is something we manage in a priority queue.
type Item struct {
	value    interface{} // The value of the item; business data.
	time_stamp int64    // The priority of the item in the queue.
	price float64	 // same as up
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}
// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return (pq[i].time_stamp < pq[j].time_stamp) && (pq[i].price > pq[j].price)
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(item *Item, value interface{}, priority int64) {
	item.value = value
	item.time_stamp = priority
	heap.Fix(pq, item.index)
}

// new queue Node
func NewQueueNode(val interface{}, tmstamp int64, price float64) *Item{
	return &Item{
		value:val,
		time_stamp:tmstamp,
		price:price,
	}
}

// create priority Queue
func NewPQ() PriorityQueue {
	return make(PriorityQueue, 32)
}
