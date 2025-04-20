package chrono

import (
    "math"
    "math/rand/v2"
    "time"
)

// StandardExponentialBackoff 提供标准指数退避算法，用于计算下一次重试的时间间隔。
//
// 参数 count 表示当前重试次数，maxRetries 指定最大重试次数，当为负数时表示无限重试。
// baseDelay 是基础延迟时间，maxDelay 是允许的最大延迟时间。
// 该函数使用默认的乘数 2 和随机化因子 0.5 来计算每次重试的具体延迟时间。
//
// 关键行为说明：
//  - 当达到最大重试次数时，返回 -1 表示不再重试
//  - 延迟时间会根据重试次数按指数增长，但不会超过 maxDelay
//  - 随机化因子引入了轻微的抖动，以避免多个任务同时触发
//
// 使用建议：
//  - 设置合理的 maxDelay 以防止过长的等待时间
//  - 对于需要快速响应的场景，可以适当减小 baseDelay
func StandardExponentialBackoff(count, maxRetries int, baseDelay, maxDelay time.Duration) time.Duration {
    return ExponentialBackoff(count, maxRetries, baseDelay, maxDelay, 2, 0.5)
}

// ExponentialBackoff 根据指数退避算法计算下一次重试的时间间隔。
//
// count 参数表示当前重试次数，maxRetries 指定最大重试次数，当为负数时表示无限重试。
// baseDelay 是基础延迟时间，maxDelay 是允许的最大延迟时间。
// multiplier 为每次重试时延迟的乘数因子，randomization 引入随机化抖动。
//
// 关键行为说明：
//  - 当达到最大重试次数时，返回 -1 表示不再重试
//  - 延迟时间会根据重试次数按指数增长，但不会超过 maxDelay
//  - 随机化因子引入了轻微的抖动，以避免多个任务同时触发
//
// 使用建议：
//  - 设置合理的 maxDelay 以防止过长的等待时间
//  - 对于需要快速响应的场景，可以适当减小 baseDelay
func ExponentialBackoff(count, maxRetries int, baseDelay, maxDelay time.Duration, multiplier, randomization float64) time.Duration {
    for {
        if count > maxRetries && maxRetries > -1 {
            return -1
        }

        delay := float64(baseDelay) * math.Pow(multiplier, float64(count))
        jitter := (rand.Float64() - 0.5) * randomization * float64(baseDelay)
        sleepDuration := time.Duration(delay + jitter)

        if sleepDuration > maxDelay {
            sleepDuration = maxDelay
        }

        return sleepDuration
    }
}
