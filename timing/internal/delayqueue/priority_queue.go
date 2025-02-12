package delayqueue

import "container/heap"

func newPriorityQueue[T any](capacity int) priorityQueue[T] {
	return make(priorityQueue[T], 0, capacity)
}

// priorityQueue 是一个最小堆实现的优先队列
type priorityQueue[T any] []*priorityQueueItem[T]

func (pq *priorityQueue[T]) Len() int {
	return len(*pq)
}

func (pq *priorityQueue[T]) Less(i, j int) bool {
	q := *pq
	return q[i].Priority < q[j].Priority
}

func (pq *priorityQueue[T]) Swap(i, j int) {
	q := *pq
	q[i], q[j] = q[j], q[i]
}

func (pq *priorityQueue[T]) Push(x interface{}) {
	n := len(*pq)
	c := cap(*pq)
	if n+1 > c {
		npq := make(priorityQueue[T], n, c*2)
		copy(npq, *pq)
		*pq = npq
	}
	*pq = (*pq)[0 : n+1]
	item := x.(*priorityQueueItem[T])
	(*pq)[n] = item
}

func (pq *priorityQueue[T]) Pop() interface{} {
	n := len(*pq)
	c := cap(*pq)
	if n < (c/2) && c > 25 {
		npq := make(priorityQueue[T], n, c/2)
		copy(npq, *pq)
		*pq = npq
	}
	item := (*pq)[n-1]
	*pq = (*pq)[0 : n-1]
	return item
}

// PeekAndShift 返回优先队列中的第一个元素，并将其从队列中移除
//   - 如果元素的优先级大于 max，则返回 nil 和优先级当前的差值
func (pq *priorityQueue[T]) PeekAndShift(max int64) (*priorityQueueItem[T], int64) {
	if pq.Len() == 0 {
		return nil, 0
	}

	item := (*pq)[0]
	if item.Priority > max {
		return item, item.Priority - max
	}
	heap.Remove(pq, 0)

	return item, 0
}
