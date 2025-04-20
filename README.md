# Chrono

Chrono 是一个全面的 Go 语言时间工具库，提供了丰富的时间操作、调度和重试机制功能。

## 特性

- **时间操作工具**：提供各种时间点计算、比较和操作的实用函数
- **高效的定时器时间轮**：用于管理大量定时任务的高性能调度器
- **Cron 表达式支持**：使用标准 cron 表达式进行任务调度
- **指数退避算法**：实现可靠的重试机制，适用于网络请求等场景

## 安装

```bash
go get github.com/kercylan98/chrono
```

## 使用示例

### 基础时间操作

```go
package main

import (
    "fmt"
    "time"

    "github.com/kercylan98/chrono"
)

func main() {
    now := time.Now()

    // 计算下一个特定时刻
    nextTime := chrono.NextMoment(now, 8, 0, 0) // 下一个早上8点
    fmt.Println("下一个早上8点:", nextTime)

    // 判断时刻是否已过
    isElapsed := chrono.Elapsed(now, 7, 0, 0)
    fmt.Println("早上7点是否已过:", isElapsed)

    // 获取时间单位的起始点
    dayStart := chrono.StartOf(now, chrono.UnitDay)
    fmt.Println("今天的开始时间:", dayStart)

    // 获取时间单位的结束点
    dayEnd := chrono.EndOf(now, chrono.UnitDay)
    fmt.Println("今天的结束时间:", dayEnd)

    // 计算两个时间点之间的差值
    delta := chrono.Delta(now, now.Add(time.Hour))
    fmt.Println("时间差:", delta)
}
```

### 使用时间轮进行任务调度

```go
package main

import (
    "fmt"
    "time"

    "github.com/kercylan98/chrono/timing"
)

func main() {
    // 创建时间轮
    wheel := timing.New()

    // 方法1：使用函数直接作为任务（无需实现Task接口）
    wheel.After(2*time.Second, timing.TaskFN(func() {
        fmt.Println("延迟任务执行于:", time.Now())
    }))

    // 方法2：使用辅助函数创建循环任务（无需实现LoopTask接口）
    loopTask := timing.NewForeverLoopTask(5*time.Second, timing.TaskFN(func() {
        fmt.Println("循环任务执行于:", time.Now())
    }))
    wheel.Loop(1*time.Second, loopTask)

    // 使用Cron表达式执行任务
    wheel.Cron("0 */5 * * * *", timing.TaskFN(func() {
        fmt.Println("定时任务执行于:", time.Now())
    }))

    // 使用命名空间组织任务
    namedWheel := wheel.Named("my-tasks")
    namedWheel.After("daily-task", 3*time.Second, timing.TaskFN(func() {
        fmt.Println("命名任务执行于:", time.Now())
    }))

    // 停止特定命名任务
    // namedWheel.Stop("daily-task")

    // 防止程序退出
    select {}
}
```

如果您仍然需要更多控制，也可以通过实现 `Task` 和 `LoopTask` 接口来创建自定义任务：

```go
package main

import (
    "fmt"
    "time"

    "github.com/kercylan98/chrono/timing"
)

type MyTask struct{}

func (t *MyTask) Execute() {
    fmt.Println("任务执行于:", time.Now())
}

type MyLoopTask struct{}

func (t *MyLoopTask) Execute() {
    fmt.Println("循环任务执行于:", time.Now())
}

func (t *MyLoopTask) Next(prev time.Time) time.Time {
    return prev.Add(5 * time.Second) // 每5秒执行一次
}

func main() {
    wheel := timing.New()
    wheel.After(2*time.Second, &MyTask{})
    wheel.Loop(1*time.Second, &MyLoopTask{})

    select {}
}
```

### 使用指数退避算法

```go
package main

import (
    "fmt"
    "time"

    "github.com/kercylan98/chrono"
)

func main() {
    // 使用标准指数退避算法
    for retryCount := 0; retryCount < 5; retryCount++ {
        // 模拟操作失败
        fmt.Println("尝试操作，次数:", retryCount+1)

        // 计算下次重试的等待时间
        delay := chrono.StandardExponentialBackoff(
            retryCount,    // 当前重试次数
            5,             // 最大重试次数
            time.Second,   // 基础延迟
            time.Minute,   // 最大延迟
        )

        if delay < 0 {
            fmt.Println("达到最大重试次数，停止重试")
            break
        }

        fmt.Printf("等待 %v 后重试\n", delay)
        time.Sleep(delay)
    }
}
```

## 高级配置

时间轮支持多种配置选项，可以根据需要进行调整：

```go
wheel := timing.GetBuilder().FromConfigurators(
    timing.WithTickDuration(100 * time.Millisecond), // 设置时间轮刻度
    timing.WithWheelSize(1024),                      // 设置时间轮大小
)
```

## 许可证

本项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE) 文件。
