package chrono

import (
	"github.com/kercylan98/chrono/src/internal/delayqueue"
	"time"
)

var (
	_                  TimingWheel = (*timingWheel)(nil)
	timingWheelBuilder             = &TimingWheelBuilder{}
)

// NewTimingWheel 创建一个用于管理大量定时任务的定时器时间轮
func NewTimingWheel(configurator ...TimingWheelConfigurator) TimingWheel {
	builder := GetTimingWheelBuilder()
	if len(configurator) > 0 {
		return builder.FromConfigurators(configurator...)
	}
	return builder.Build()
}

// GetTimingWheelBuilder 获取一个用于创建时间轮的构建器
func GetTimingWheelBuilder() *TimingWheelBuilder {
	return timingWheelBuilder
}

// TimingWheelBuilder NewTimingWheel 创建一个用于管理大量定时任务的定时器时间轮的构建器
type TimingWheelBuilder struct{}

// Build 创建一个默认配置的时间轮
func (builder *TimingWheelBuilder) Build() TimingWheel {
	tw := &timingWheel{}
	tw.timingWheelInternal = newTimingWheelInternal(tw, NewTimingWheelConfig())
	tw.init(0, nil)
	return tw
}

// build 内部构建方法
func (builder *TimingWheelBuilder) build(startMs int64, queue *delayqueue.DelayQueue[bucket], configuration TimingWheelConfiguration) TimingWheel {
	tw := &timingWheel{}
	tw.timingWheelInternal = newTimingWheelInternal(tw, configuration)
	tw.init(startMs, queue)
	return tw
}

// FromConfiguration 从配置中创建一个时间轮
func (builder *TimingWheelBuilder) FromConfiguration(config TimingWheelConfiguration) TimingWheel {
	tw := &timingWheel{}
	tw.timingWheelInternal = newTimingWheelInternal(tw, config)
	tw.init(0, nil)
	return tw
}

// FromCustomize 通过自定义配置构建时间轮
func (builder *TimingWheelBuilder) FromCustomize(configuration TimingWheelConfiguration, configurators ...TimingWheelConfigurator) TimingWheel {
	for _, configurator := range configurators {
		configurator.Configure(configuration)
	}
	return builder.FromConfiguration(configuration)
}

// FromConfigurators 从配置器中创建一个时间轮
func (builder *TimingWheelBuilder) FromConfigurators(configurators ...TimingWheelConfigurator) TimingWheel {
	var config = NewTimingWheelConfig()
	for _, c := range configurators {
		c.Configure(config)
	}
	return builder.FromConfiguration(config)
}

// TimingWheel 用于管理大量定时任务的定时器时间轮，它是一个时间轮的抽象
type TimingWheel interface {
	timingWheelInternal

	// AfterFunc 创建一个在一段时间后执行的任务
	AfterFunc(duration time.Duration, task func()) Timer
}

// timingWheel 是 TimingWheel 的默认实现
type timingWheel struct {
	timingWheelInternal
}

func (t *timingWheel) AfterFunc(duration time.Duration, task func()) Timer {
	timer := newTimer(ToMillisecond(time.Now().Add(duration)), task)
	t.contract(timer)
	return timer
}
