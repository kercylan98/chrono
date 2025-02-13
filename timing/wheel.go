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

	// After 创建一个在一段时间后执行的任务
	After(duration time.Duration, task Task) Timer

	// Loop 创建一个循环执行的任务，它将在 duration 时间后首次执行，然后根据 LoopTask.Next 方法返回的时间再次执行
	Loop(duration time.Duration, task LoopTask) Timer

	// Cron 通过 cron 表达式创建一个任务，当表达式无效时将返回错误
	//  - 表达式说明可参阅：https://github.com/gorhill/cronexpr
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
	timer := newTimer(chrono.TimeToMillisecond(time.Now().Add(duration)), task.Execute)
	t.contract(timer)
	return timer
}

func (t *wheel) Loop(duration time.Duration, task LoopTask) Timer {
	var timer Timer
	timer = newTimer(chrono.TimeToMillisecond(time.Now().Add(duration)), func() {
		defer func() {
			previous := chrono.MillisecondToTime(timer.getExpiration())
			next := task.Next(previous)
			if !next.IsZero() && next.After(previous) {
				timer.setExpiration(chrono.TimeToMillisecond(next))
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
	timer = newTimer(chrono.TimeToMillisecond(expression.Next(now)), func() {
		defer func() {
			next := expression.Next(now)
			timer.setExpiration(chrono.TimeToMillisecond(next))
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
