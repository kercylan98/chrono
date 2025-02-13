package timing

import "time"

// Task 是一个任务，它被用来在计时器到达指定的过期时间时执行
type Task interface {
	// Execute 执行任务
	Execute()
}

type TaskFn func()

func (f TaskFn) Execute() {
	f()
}

// LoopTask 是一个循环任务，它被用来在计时器到达指定的过期时间时执行，并且可以指定下一次执行的时间
type LoopTask interface {
	Task

	// Next 返回下一次执行的时间
	//  - 参数 previous 表示了上一次的执行时间，当返回的时间小于 previous 时，任务将不再执行
	Next(previous time.Time) time.Time
}

// NewLoopTask 创建一个循环任务，它将根据 interval 作为间隔时间反复执行最多 times 次
//   - 当 times 的值为 0 时，任务将不会被执行，如果 times 的值为负数，那么任务将永远不会停止，除非主动停止
func NewLoopTask(interval time.Duration, times int, task Task) LoopTask {
	return &loopTask{
		interval: interval,
		times:    times,
		task:     task,
	}
}

// NewForeverLoopTask 创建一个永久循环任务，它将根据 interval 作为间隔时间无限循环执行
func NewForeverLoopTask(interval time.Duration, task Task) LoopTask {
	return NewLoopTask(interval, -1, task)
}

type loopTask struct {
	interval time.Duration
	times    int
	task     Task
}

func (f *loopTask) Next(previous time.Time) time.Time {
	if f.times == 0 {
		return time.Time{}
	}
	return previous.Add(f.interval)
}

func (f *loopTask) Execute() {
	if f.times == 0 {
		return
	}
	f.task.Execute()
	if f.times > 0 {
		f.times--
	}
}
