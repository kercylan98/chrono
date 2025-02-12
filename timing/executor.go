package timing

type Executor interface {
	// Execute 执行任务
	Execute(task func())
}

type ExecutorFn func(task func())

func (f ExecutorFn) Execute(task func()) {
	f(task)
}
