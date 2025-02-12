package chrono

type TimingWheelExecutor interface {
	// Execute 执行任务
	Execute(task func())
}

type TimingWheelExecutorFn func(task func())

func (f TimingWheelExecutorFn) Execute(task func()) {
	f(task)
}
