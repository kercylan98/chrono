package timing

import (
	"sync"
	"time"
)

type Named interface {
	wheelInternal

	// After 创建一个在一段时间后执行的任务
	After(name string, duration time.Duration, task Task)

	// Loop 创建一个循环执行的任务，它将在 duration 时间后首次执行，然后根据 LoopTask.Next 方法返回的时间再次执行
	Loop(name string, duration time.Duration, task LoopTask)

	// Cron 通过 cron 表达式创建一个任务，当表达式无效时将返回错误
	//  - 表达式说明可参阅：https://github.com/gorhill/cronexpr
	Cron(name string, cron string, task Task) error

	// Stop 停止指定名称的任务
	Stop(name string)

	// Clear 清除所有任务
	Clear()

	// Timer 获取使用 Timer 维护任务的时间轮 API
	Timer() Wheel
}

func newNamed(t Wheel) Named {
	return &named{
		Wheel:  t,
		timers: make(map[string]Timer),
	}
}

type named struct {
	Wheel
	timers map[string]Timer
	lock   sync.RWMutex
}

func (t *named) After(name string, duration time.Duration, task Task) {
	t.lock.Lock()
	if old, ok := t.timers[name]; ok {
		old.Stop()
	}
	t.timers[name] = t.Wheel.After(duration, task)
	t.lock.Unlock()
}

func (t *named) Loop(name string, duration time.Duration, task LoopTask) {
	t.lock.Lock()
	if old, ok := t.timers[name]; ok {
		old.Stop()
	}
	t.timers[name] = t.Wheel.Loop(duration, task)
	t.lock.Unlock()
}

func (t *named) Cron(name string, cron string, task Task) error {
	if timer, err := t.Wheel.Cron(cron, task); err != nil {
		return err
	} else {
		t.lock.Lock()
		if old, ok := t.timers[name]; ok {
			old.Stop()
		}
		t.timers[name] = timer
		t.lock.Unlock()
	}
	return nil
}

func (t *named) Stop(name string) {
	t.lock.Lock()
	if timer, ok := t.timers[name]; ok {
		timer.Stop()
	}
	delete(t.timers, name)
	t.lock.Unlock()
}

func (t *named) Clear() {
	t.lock.Lock()
	for _, timer := range t.timers {
		timer.Stop()
	}
	t.timers = make(map[string]Timer)
	t.lock.Unlock()
}

func (t *named) Timer() Wheel {
	return t.Wheel
}
