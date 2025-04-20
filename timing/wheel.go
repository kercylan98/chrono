package timing

import (
    "github.com/gorhill/cronexpr"
    "github.com/kercylan98/chrono"
    "github.com/kercylan98/chrono/timing/internal/delayqueue"
    "sync"
    "time"
)

var (
    _       Wheel = (*wheel)(nil)
    builder       = &Builder{}
)

// New 创建一个用于管理大量定时任务的定时器时间轮
func New(configurator ...Configurator) Wheel {
    builder := GetBuilder()
    if len(configurator) > 0 {
        return builder.FromConfigurators(configurator...)
    }
    return builder.Build()
}

// GetBuilder 获取一个用于创建时间轮的构建器
func GetBuilder() *Builder {
    return builder
}

// Builder New 创建一个用于管理大量定时任务的定时器时间轮的构建器
type Builder struct{}

// Build 创建一个默认配置的时间轮
func (builder *Builder) Build() Wheel {
    tw := &wheel{}
    tw.wheelInternal = newWheelInternal(tw, NewConfig())
    tw.init(0, nil)
    return tw
}

// build 内部构建方法
func (builder *Builder) build(startMs int64, queue *delayqueue.DelayQueue[bucket], configuration Configuration) Wheel {
    tw := &wheel{}
    tw.wheelInternal = newWheelInternal(tw, configuration)
    tw.init(startMs, queue)
    return tw
}

// FromConfiguration 从配置中创建一个时间轮
func (builder *Builder) FromConfiguration(config Configuration) Wheel {
    tw := &wheel{}
    tw.wheelInternal = newWheelInternal(tw, config)
    tw.init(0, nil)
    return tw
}

// FromCustomize 通过自定义配置构建时间轮
func (builder *Builder) FromCustomize(configuration Configuration, configurators ...Configurator) Wheel {
    for _, configurator := range configurators {
        configurator.Configure(configuration)
    }
    return builder.FromConfiguration(configuration)
}

// FromConfigurators 从配置器中创建一个时间轮
func (builder *Builder) FromConfigurators(configurators ...Configurator) Wheel {
    var config = NewConfig()
    for _, c := range configurators {
        c.Configure(config)
    }
    return builder.FromConfiguration(config)
}

// Wheel 用于管理大量定时任务的定时器时间轮，它是一个时间轮的抽象
type Wheel interface {
    wheelInternal

    // After 创建一个在指定延迟后执行的任务。
    //
    // duration 参数定义了任务首次执行前的等待时间，若为零或负值则立即执行。
    // 任务通过 Task 接口定义，Execute 方法将在延迟结束后被调用。
    // 返回 Timer 对象用于控制任务状态，如停止或检查是否已停止。
    //
    // 关键行为说明：
    //  - 若 duration 为零或负值，任务将立即执行
    //  - 使用返回的 Timer 可以停止任务
    //  - 任务执行过程中发生 panic 将被捕获并记录，但不会中断调度
    After(duration time.Duration, task Task) Timer

    // Loop 创建并启动一个循环任务，根据指定的初始延迟和任务定义执行。
    //
    // duration 参数指定了首次执行前的等待时间，设置为零或负值将立即触发执行。
    // task 参数是一个实现了 LoopTask 接口的任务，定义了任务的具体行为及下次执行的时间。
    //
    // 关键行为说明：
    //  - 当 duration <= 0 时，任务将立即执行
    //  - 使用返回的 Timer 可以停止任务
    //  - 异常处理机制会捕获执行过程中的 panic 并记录，但不影响后续调度
    Loop(duration time.Duration, task LoopTask) Timer

    // Cron 通过 cron 表达式创建一个周期性任务。
    //
    // 参数 cron 是一个标准的 cron 表达式，用于定义任务的执行时间。task 参数是实际执行的任务。
    // 如果 cron 表达式无效，将返回错误。
    //
    // 时间参数精度取决于系统时钟，实际执行可能存在毫秒级偏差。
    Cron(cron string, task Task) (Timer, error)

    // Named 获取使用命名维护任务的时间轮 API
    //   - 当 topic 不为空时，将返回一个命名空间为 topic 的 Named 实例，不同的 Named 实例之间的任务不会相互影响
    Named(topic ...string) Named
}

// wheel 是 Wheel 的默认实现
type wheel struct {
    wheelInternal
    named map[string]Named
    rw    sync.RWMutex
}

func (t *wheel) After(duration time.Duration, task Task) Timer {
    timer := newTimer(chrono.ToMillisecond(time.Now().Add(duration)), task.Execute)
    t.contract(timer)
    return timer
}

func (t *wheel) Loop(duration time.Duration, task LoopTask) Timer {
    var timer Timer
    timer = newTimer(chrono.ToMillisecond(time.Now().Add(duration)), func() {
        defer func() {
            previous := chrono.ToTime(timer.getExpiration())
            next := task.Next(previous)
            if !next.IsZero() && next.After(previous) {
                timer.setExpiration(chrono.ToMillisecond(next))
                t.contract(timer)
            }
        }()

        task.Execute()
    })
    t.contract(timer)
    return timer
}

func (t *wheel) Cron(cron string, task Task) (Timer, error) {
    expression, err := cronexpr.Parse(cron)
    if err != nil {
        return nil, err
    }
    var now = time.Now()
    var timer Timer
    timer = newTimer(chrono.ToMillisecond(expression.Next(now)), func() {
        defer func() {
            next := expression.Next(now)
            timer.setExpiration(chrono.ToMillisecond(next))
            t.contract(timer)
        }()

        task.Execute()
    })
    t.contract(timer)
    return timer, nil
}

func (t *wheel) Named(topic ...string) Named {
    t.rw.Lock()
    defer t.rw.Unlock()
    var name string
    if len(topic) > 0 {
        name = topic[0]
    }
    if t.named == nil {
        t.named = make(map[string]Named)
    }

    if named, exist := t.named[name]; exist {
        return named
    } else {
        named = newNamed(t)
        t.named[name] = named
        return named
    }
}
