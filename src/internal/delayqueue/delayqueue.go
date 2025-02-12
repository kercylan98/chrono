package delayqueue

import (
	"container/heap"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const (
	delayQueueSleeping = iota
	delayQueueWorking
)

func New[T QueueItem](size int, timeGetter func() int64, handler func(v T)) *DelayQueue[T] {
	return &DelayQueue[T]{
		priorityQueue: newPriorityQueue[T](size),
		timeGetter:    timeGetter,
		handler:       handler,
		wakeupCancel:  func() {},
	}
}

type QueueItem interface {
	Size() int
}

type DelayQueue[T QueueItem] struct {
	state         atomic.Int32
	n             atomic.Int64
	mu            sync.Mutex
	priorityQueue priorityQueue[T]
	timeGetter    func() int64
	handler       func(v T)
	wakeupCtx     context.Context
	wakeupCancel  context.CancelFunc
}

// Add 将元素插入到当前队列中。
func (q *DelayQueue[T]) Add(elem T, expiration int64) {
	item := newPriorityQueueItem(elem, expiration)

	q.mu.Lock()
	heap.Push(&q.priorityQueue, item)
	q.mu.Unlock()

	if q.state.CompareAndSwap(delayQueueSleeping, delayQueueWorking) {
		go q.wakeup()
	} else {
		q.n.Add(1)
		q.wakeupCancel()
	}
}

// Refresh 刷新元素的过期时间。
func (q *DelayQueue[T]) Refresh() {
	q.wakeupCancel()
}

func (q *DelayQueue[T]) wakeup() {
	for {
		q.process()
		q.state.Store(delayQueueSleeping)
		if q.n.Load() == 0 {
			break
		} else if !q.state.CompareAndSwap(delayQueueSleeping, delayQueueWorking) {
			break
		}
	}
}

func (q *DelayQueue[T]) process() {
	q.n.Store(0)

	for {
		now := q.timeGetter()

		q.mu.Lock()
		item, delta := q.priorityQueue.PeekAndShift(now)
		q.mu.Unlock()

		if item == nil || item.Value.Size() == 0 {
			break // 没有任何元素待处理
		}

		if delta > 0 {

			after := time.Duration(delta)
			q.wakeupCtx, q.wakeupCancel = context.WithTimeout(context.Background(), after)
			select {
			case <-q.wakeupCtx.Done():
				continue
			}
		}

		q.handler(item.Value)

	}
}
