package chrono

import (
	"container/list"
	"github.com/kercylan98/chrono/src/internal/delayqueue"
	"sync"
	"sync/atomic"
)

var (
	_ bucket = (*bucketImpl)(nil)
)

func newBucket(wheel TimingWheel) bucket {
	return &bucketImpl{
		wheel:  wheel,
		timers: list.New(),
	}
}

// bucket 计时桶是一个计时器的集合，它管理了一组相同过期时间的计时器
type bucket interface {
	delayqueue.QueueItem

	// getExpiration 返回计时桶的毫秒级过期时间
	getExpiration() int64

	// setExpiration 设置计时桶的毫秒级过期时间，当过期时间发生变化时返回 true
	setExpiration(expiration int64) bool

	// add 添加一个计时器到计时桶中
	add(timer Timer)

	// unlockedAdd 添加一个计时器到计时桶中，该函数不会加锁
	unlockedAdd(timer Timer)

	// remove 从计时桶中移除一个计时器，如果计时器不在计时桶中则返回 false
	remove(Timer) bool

	// flush 清空计时桶中的所有计时器，并将这些计时器重新插入到时间轮中
	flush(adder func(Timer))
}

type bucketImpl struct {
	expiration atomic.Int64
	timers     *list.List
	rw         sync.RWMutex
	wheel      TimingWheel // 所属时间轮
}

func (b *bucketImpl) Size() int {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.timers.Len()
}

func (b *bucketImpl) getExpiration() int64 {
	return b.expiration.Load()
}

func (b *bucketImpl) setExpiration(expiration int64) bool {
	return b.expiration.Swap(expiration) != expiration
}

func (b *bucketImpl) add(timer Timer) {
	b.rw.Lock()
	e := b.timers.PushBack(timer)
	b.rw.Unlock()

	timer.setBucket(b, e)
}

func (b *bucketImpl) unlockedAdd(timer Timer) {
	e := b.timers.PushBack(timer)
	timer.setBucket(b, e)
}

func (b *bucketImpl) unlockedRemove(t Timer) bool {
	if t.getBucket() != b {
		return false
	}

	b.timers.Remove(t.getElement())

	t.setBucket(nil, nil)
	b.wheel.refreshDelayQueue()
	return true
}

func (b *bucketImpl) remove(t Timer) bool {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.unlockedRemove(t)
}

func (b *bucketImpl) flush(adder func(Timer)) {
	// 该函数会在延迟队列的回调中被调用，该调用是异步的，需要确保线程安全
	b.rw.Lock()
	defer b.rw.Unlock()

	for e := b.timers.Front(); e != nil; {
		next := e.Next()

		t := e.Value.(Timer)
		b.unlockedRemove(t)

		// 添加到时间轮中时，如果任务时间已经到达，将被执行
		adder(t)

		e = next
	}

	b.setExpiration(-1)
}
