package chrono

import (
	"github.com/kercylan98/options"
	"time"
)

var (
	_                          TimingWheelConfiguration = (*timingWheelConfiguration)(nil)
	defaultTimingWheelExecutor                          = TimingWheelExecutorFn(func(task func()) {
		task()
	})
)

// NewTimingWheelConfig 创建一个用于 TimingWheel 的默认配置器
func NewTimingWheelConfig() TimingWheelConfiguration {
	c := &timingWheelConfiguration{
		tick:     1,
		size:     20,
		executor: defaultTimingWheelExecutor,
	}
	c.LogicOptions = options.NewLogicOptions[TimingWheelOptionsFetcher, TimingWheelOptions](c, c)
	return c
}

// TimingWheelConfigurator 是 TimingWheel 的配置接口，它允许结构化的配置 TimingWheel
type TimingWheelConfigurator interface {
	// Configure 配置 TimingWheel
	Configure(config TimingWheelConfiguration)
}

// TimingWheelConfiguratorFn 是 TimingWheel 的配置接口，它允许通过函数式的方式配置 TimingWheel
type TimingWheelConfiguratorFn func(config TimingWheelConfiguration)

func (f TimingWheelConfiguratorFn) Configure(config TimingWheelConfiguration) {
	f(config)
}

type TimingWheelConfiguration interface {
	TimingWheelOptions
	TimingWheelOptionsFetcher
}

type TimingWheelOptions interface {
	options.LogicOptions[TimingWheelOptionsFetcher, TimingWheelOptions]

	// WithTick 设置时间轮的刻度，单位为毫秒
	WithTick(tick time.Duration) TimingWheelConfiguration

	// WithSize 设置时间轮的大小
	WithSize(size int) TimingWheelConfiguration

	// WithExecutor 设置时间轮的执行器
	WithExecutor(executor TimingWheelExecutor) TimingWheelConfiguration
}

type TimingWheelOptionsFetcher interface {
	FetchTick() int64

	FetchSize() int64

	FetchExecutor() TimingWheelExecutor
}

type timingWheelConfiguration struct {
	options.LogicOptions[TimingWheelOptionsFetcher, TimingWheelOptions]
	tick     int64 // 每个刻度的毫秒级时间
	size     int64 // 每个时间轮的毫秒级间隔时间
	executor TimingWheelExecutor
}

func (t *timingWheelConfiguration) WithTick(tick time.Duration) TimingWheelConfiguration {
	t.tick = int64(tick / time.Millisecond)
	return t
}

func (t *timingWheelConfiguration) WithSize(size int) TimingWheelConfiguration {
	t.size = int64(size)
	return t
}

func (t *timingWheelConfiguration) WithExecutor(executor TimingWheelExecutor) TimingWheelConfiguration {
	t.executor = executor
	return t
}

func (t *timingWheelConfiguration) FetchTick() int64 {
	return t.tick
}

func (t *timingWheelConfiguration) FetchSize() int64 {
	return t.size
}

func (t *timingWheelConfiguration) FetchExecutor() TimingWheelExecutor {
	return t.executor
}
