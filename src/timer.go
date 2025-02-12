package chrono

import (
	"container/list"
	"sync/atomic"
)

// Timer 是一个计时器，它可以在到达指定的过期时间时触发一个事件
type Timer interface {
	// Stop 停止计时器，如果计时器已经停止则返回 false
	Stop() bool

	getExpiration() int64

	getTask() func()

	getBucket() bucket

	getElement() *list.Element

	setBucket(bucket bucket, element *list.Element)
}

func newTimer(expiration int64, task func()) Timer {
	return &timerImpl{
		expiration: expiration,
		task:       task,
	}
}

type timerImpl struct {
	expiration int64                  // 过期时间
	task       func()                 // 任务
	bucket     atomic.Pointer[bucket] // 所在的桶
	element    *list.Element          // 桶元素
}

func (t *timerImpl) getExpiration() int64 {
	return t.expiration
}

func (t *timerImpl) Stop() bool {
	stopped := false
	for b := t.getBucket(); b != nil; b = t.getBucket() {
		stopped = b.remove(t)
	}
	return stopped
}

func (t *timerImpl) getTask() func() {
	return t.task
}

func (t *timerImpl) getBucket() bucket {
	b := t.bucket.Load()
	if b == nil {
		return nil
	}
	return *b
}

func (t *timerImpl) setBucket(bucket bucket, element *list.Element) {
	t.bucket.Store(&bucket)
	t.element = element
}

func (t *timerImpl) getElement() *list.Element {
	return t.element
}
