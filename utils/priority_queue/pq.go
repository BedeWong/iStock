package priority_queue

import (
	"container/heap"
	"math"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/pkg/errors"
	"fmt"
)

// 根据 官方文档改编的 优先队列：
// An Item is something we manage in a priority queue.
type Item struct {
	value    interface{} // The value of the item; business data.
	time_stamp int64    // The priority of the item in the queue.
	price float64	 // same as up
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// 獲取 value 字段
func (it Item)Value()interface{} {
	return it.value
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

const (
	price_min = 0.0
	price_max = math.MaxFloat64

	time_stamp_min = 0
	time_stamp_max = math.MaxInt64
)

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
func (pq *PriorityQueue) Update(item *Item, value interface{}, price float64, priority int64) {
	item.value = value
	item.price = price
	item.time_stamp = priority
	heap.Fix(pq, item.index)
}

// 删除一个元素，
func (pq *PriorityQueue) Remove(cmp func(val interface{})bool) error{
	for _, it := range *pq {
		log.Debug("pd Remove: range item: %#v", it)

		if cmp(it.value) {
			pq.Update(it, it.value, price_min, time_stamp_min)  // 修改 优先级到 top1

			el := pq.Pop()
			item, _ := el.(*Item)

			if item.time_stamp != time_stamp_min || item.price != price_min {
				log.Error("元素不在pop()中。")
				return errors.New(fmt.Sprintf("Remove:执行失败。"))
			}
			// 结束函数
			return nil
		}
	}

	return errors.New("item not found.")
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
	return make(PriorityQueue, 0)
}
