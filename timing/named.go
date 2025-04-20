package timing

import (
    "sync"
    "time"
)

// Named 提供了命名任务的管理接口，支持定时、循环和基于cron表达式的任务调度。
//
// 通过 After 方法可以创建一个在指定时长后执行的任务。
// Loop 方法允许创建循环任务，该任务将在首次延迟后开始，并根据 LoopTask.Next 方法返回的时间间隔重复执行。
// Cron 方法则利用 cron 表达式定义更复杂的调度模式，当表达式无效时会返回错误。
// Stop 和 Clear 分别用于停止特定名称的任务或清除所有任务。
// Timer 方法提供了访问底层时间轮 API 的方式，以实现更精细的任务控制。
//
// 关键行为说明：
//  - 同名任务会被新任务覆盖，确保任务唯一性
//  - 停止任务时，正在执行的任务将完成当前操作后再退出
//  - 使用 Cron 时需保证表达式正确，否则任务不会被创建
type Named interface {
    wheelInternal

    // After 在指定时长后执行给定任务，支持同名任务覆盖。
    //
    // duration 参数定义了任务首次执行前的等待时间，零或负值表示立即执行。
    // 同名任务将被新创建的任务替换，并继承其调度队列位置。
    // 任务的实际执行可能受系统时钟精度影响而存在毫秒级偏差。
    //
    // 关键行为说明：
    //  - 父级上下文关闭时，正在执行的任务会完成当前操作后再退出
    //  - 使用建议：即时任务可设置 duration = 0 实现单次触发
    After(name string, duration time.Duration, task Task)

    // Loop 创建一个具有指定延迟和循环间隔的任务，支持同名任务覆盖。
    //
    // name 参数用于标识任务，同名任务将被新任务覆盖。duration 参数设置首次执行前的等待时间，
    // 为零或负值时立即执行。task 参数是一个 LoopTask 类型，定义了任务的具体执行逻辑及下次执行时间。
    //
    // 关键行为说明：
    //  - 同名任务会被新任务覆盖，确保唯一性
    //  - 当 duration 为零或负值时，任务会立即开始执行
    //  - 任务的执行依赖于 LoopTask.Next 方法返回的时间间隔
    //
    // 使用建议：
    //  - 确保 LoopTask 实现正确处理并发情况
    //  - 对长时间运行的任务使用上下文控制超时和取消
    Loop(name string, duration time.Duration, task LoopTask)

    // Cron 使用 cron 表达式创建具有复杂调度模式的任务。
    //
    // 参数 name 用于唯一标识任务，相同名称的任务将被新任务覆盖。
    // 参数 cron 是一个标准的 cron 表达式，定义了任务的执行时间表。
    // 参数 task 是要执行的任务，必须实现 Task 接口。
    //
    // 关键行为说明：
    //  - 同名任务会被新任务覆盖
    //  - 当父级上下文关闭时，已进入执行阶段的任务会完成当前操作再退出
    //  - 异常处理机制会捕获并记录执行过程中的 panic，但不会中断任务调度流程
    Cron(name string, cron string, task Task) error

    // Stop 停止指定名称的任务。
    //
    // name 参数用于标识要停止的任务。如果任务正在执行，它将完成当前操作后再退出。
    //
    // 关键行为说明：
    //  - 正在执行的任务会完成当前操作再退出
    Stop(name string)

    // Clear 清除所有已注册的任务。
    //
    // 该方法会立即停止并清除当前命名空间下的所有任务，包括正在执行的任务也会被取消。
    // 清除过程中不会等待任何任务完成其当前操作。
    //
    // 关键行为说明：
    //  - 正在执行的任务会完成当前操作再退出
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
