package timing

import "time"

// Task 定义了任务执行的基本接口。
//
// 该接口主要用于定义一个可以被执行的任务单元。Execute 方法用于触发任务的实际执行逻辑。
//
// 关键行为说明：
//  - Execute 方法应包含任务的具体实现逻辑
//  - 实现类需确保方法的线程安全，特别是在并发环境中
//
// 使用建议：
//  - 保持 Execute 方法轻量且高效以支持高频率调用
//  - 在复杂任务中考虑使用上下文控制超时和取消
type Task interface {
    // Execute 执行任务
    Execute()
}

// TaskFN 定义了一个无参数、无返回值的任务函数类型。
//
// 该类型用于封装具体的任务逻辑，以便在定时任务调度器中执行。通过实现 Task 接口的 Execute 方法，
// 可以将 TaskFN 类型的任务函数纳入调度流程。TaskFN 的定义非常简单，仅包含一个函数签名，
// 适用于需要周期性或一次性执行的轻量级任务场景。
//
// 关键行为说明：
//  - 任务执行时，不会传递任何参数，也不会有返回值
//  - 任务执行过程中抛出的 panic 会被捕获并记录，但不会中断任务调度
type TaskFN func()

func (f TaskFN) Execute() {
    f()
}

// LoopTask 是一个循环任务，它被用来在计时器到达指定的过期时间时执行，并且可以指定下一次执行的时间
type LoopTask interface {
    Task

    // Next 返回下一次执行的时间
    //  - 参数 previous 表示了上一次的执行时间，当返回的时间小于 previous 时，任务将不再执行
    Next(previous time.Time) time.Time
}

// NewLoopTask 创建具有生命周期管理的延迟执行任务，支持动态策略配置和同名任务替换。
//
// 任务调度策略通过参数组合实现灵活控制：interval 参数控制任务的循环间隔，当该值小于等于 0 时则任务将尽可能快地连续执行。
// times 参数限制最大执行次数，非正值时任务将持续运行直至主动终止，为零时任务将不被执行。
// task 参数指定具体要执行的任务。
//
// 时间参数精度取决于系统时钟，实际执行可能存在毫秒级偏差。
//
// 关键行为说明：
//  - 当父级上下文关闭时，已进入执行阶段的任务会完成当前操作再退出
//  - 连续执行模式中，若任务耗时超过间隔时长，下次执行将顺延至当前操作完成
func NewLoopTask(interval time.Duration, times int, task Task) LoopTask {
    return &loopTask{
        interval: interval,
        times:    times,
        task:     task,
    }
}

// NewForeverLoopTask 创建一个无限循环执行的任务，基于给定的时间间隔和任务。
//
// 该函数接受两个参数：interval 和 task。interval 参数定义了每次任务执行之间的等待时间，
// 若 interval 小于等于 0，则任务将尽可能快地连续执行。task 参数指定了要执行的具体任务。
//
// 同名任务管理采用覆盖机制，新建任务会自动取消前序同名实例并继承其调度队列位置。
// 时间参数精度取决于系统时钟，实际执行可能存在毫秒级偏差。
//
// 关键行为说明：
//  - 当父级上下文关闭时，已进入执行阶段的任务会完成当前操作再退出
//  - 连续执行模式中，若任务耗时超过间隔时长，下次执行将顺延至当前操作完成
//  - 异常处理机制会捕获并记录执行过程中的 panic，但不会中断任务调度流程
//
// 使用建议：
//  - 对于需要快速响应的场景，可以设置 interval 为负值以实现最小延迟执行
//  - 长期运行的任务应通过 context.WithTimeout 创建有界上下文来控制生命周期
//
// 并发机制采用分级协程池管理，任务提交与执行分离保障调度稳定性。
// 高频任务建议配置执行限速策略避免协程数量激增。
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
    if now := time.Now(); previous.Before(now) {
        previous = now
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
