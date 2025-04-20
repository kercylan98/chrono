package timing

import (
    "fmt"
    "runtime/debug"
)

type Executor interface {
    // Execute 执行任务
    Execute(task func())
}

type ExecutorFN func(task func())

func (f ExecutorFN) Execute(task func()) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Println(err)
            debug.PrintStack()
        }
    }()
    f(task)
}
