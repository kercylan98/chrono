package timing

import (
	"github.com/kercylan98/chrono/timing/internal/delayqueue"
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

	// AfterFunc 创建一个在一段时间后执行的任务
	AfterFunc(duration time.Duration, task func()) Timer
}

// wheel 是 Wheel 的默认实现
type wheel struct {
	wheelInternal
}

func (t *wheel) AfterFunc(duration time.Duration, task func()) Timer {
	timer := newTimer(ToMillisecond(time.Now().Add(duration)), task)
	t.contract(timer)
	return timer
}
