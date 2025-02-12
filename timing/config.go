package timing

import (
	"github.com/kercylan98/options"
	"time"
)

var (
	_               Configuration = (*configuration)(nil)
	defaultExecutor               = ExecutorFn(func(task func()) {
		task()
	})
)

// NewConfig 创建一个用于 Wheel 的默认配置器
func NewConfig() Configuration {
	c := &configuration{
		tick:     1,
		size:     20,
		executor: defaultExecutor,
	}
	c.LogicOptions = options.NewLogicOptions[OptionsFetcher, Options](c, c)
	return c
}

// Configurator 是 Wheel 的配置接口，它允许结构化的配置 Wheel
type Configurator interface {
	// Configure 配置 Wheel
	Configure(config Configuration)
}

// ConfiguratorFn 是 Wheel 的配置接口，它允许通过函数式的方式配置 Wheel
type ConfiguratorFn func(config Configuration)

func (f ConfiguratorFn) Configure(config Configuration) {
	f(config)
}

type Configuration interface {
	Options
	OptionsFetcher
}

type Options interface {
	options.LogicOptions[OptionsFetcher, Options]

	// WithTick 设置时间轮的刻度，单位为毫秒
	WithTick(tick time.Duration) Configuration

	// withTick 内部设置时间轮的刻度，单位为毫秒。该函数不进行换算
	withTick(tick int64) Configuration

	// WithSize 设置时间轮的大小
	WithSize(size int) Configuration

	// WithExecutor 设置时间轮的执行器
	WithExecutor(executor Executor) Configuration
}

type OptionsFetcher interface {
	FetchTick() int64

	FetchSize() int64

	FetchExecutor() Executor
}

type configuration struct {
	options.LogicOptions[OptionsFetcher, Options]
	tick     int64 // 每个刻度的毫秒级时间
	size     int64 // 每个时间轮的毫秒级间隔时间
	executor Executor
}

func (t *configuration) WithTick(tick time.Duration) Configuration {
	t.tick = int64(tick / time.Millisecond)
	return t
}

func (t *configuration) withTick(tick int64) Configuration {
	t.tick = tick
	return t
}

func (t *configuration) WithSize(size int) Configuration {
	t.size = int64(size)
	return t
}

func (t *configuration) WithExecutor(executor Executor) Configuration {
	t.executor = executor
	return t
}

func (t *configuration) FetchTick() int64 {
	return t.tick
}

func (t *configuration) FetchSize() int64 {
	return t.size
}

func (t *configuration) FetchExecutor() Executor {
	return t.executor
}
