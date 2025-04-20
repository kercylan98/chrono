package chrono

import (
    "time"
)

// NewPeriod 创建一个新的时间段，确保开始时间不晚于结束时间。
//
// 如果给定的开始时间晚于结束时间，则自动交换两者以保证正确的顺序。
// 该函数用于构建一个合法的时间段对象，适用于需要明确时间范围的各种场景。
//
// 关键行为说明：
//  - 自动调整输入参数顺序，确保始终返回有效的时间段
//
// 使用建议：
// 确保传入的时间点是基于同一时区和日历系统，避免因时区差异导致逻辑错误。
func NewPeriod(start, end time.Time) Period {
    if start.After(end) {
        start, end = end, start
    }
    return Period{start, end}
}

// Period 表示一个时间区间，由开始时间和结束时间组成。
//
// 时间区间的开始和结束时间通过两个 time.Time 类型的值表示。
// 如果创建时开始时间晚于结束时间，两者会自动交换以保证正确的顺序。
// 使用此类型可以方便地进行时间段内的各种计算和判断。
//
// 关键行为说明：
//  - 时间区间的持续时间可通过 Duration 方法获取
//  - 可通过 IsZero 判断是否为零值
//  - 通过 IsInvalid 判断是否包含无效的时间
//  - 提供多种方法来判断时间段与指定时间的关系，如 IsBefore, IsAfter, IsBetween 等
//
// 使用建议：
//  - 确保时间区间不为空或无效，以避免在后续操作中出现错误
//  - 在处理时间相关的逻辑时，优先使用该类型的内置方法，以确保准确性和一致性
//
// 并发机制方面，由于是简单的数据结构，通常不需要特别的并发控制。
type Period [2]time.Time

// Start 返回时间段的开始时间。
//
// 该方法直接返回 Period 结构体中的第一个 time.Time 值，表示时间段的起始点。
// 如果在创建 Period 时，开始时间晚于结束时间，两者会自动交换以确保正确的顺序。
//
// 关键行为说明：
//  - 调用此方法不会改变 Period 的内部状态
//  - 若需要获取整个时间段的信息，请结合使用 End 方法
//
// 使用建议：
//  - 确保 Period 实例有效且非零值，以避免返回无效的时间
//  - 在处理时间相关的逻辑时，优先使用内置方法以确保准确性和一致性
func (p Period) Start() time.Time {
    return p[0]
}

// End 返回时间段的结束时间。
//
// 该方法直接返回 Period 结构体中的第二个 time.Time 值，表示时间段的终点。
// 如果在创建 Period 时，开始时间晚于结束时间，两者会自动交换以确保正确的顺序。
//
// 关键行为说明：
//  - 调用此方法不会改变 Period 的内部状态
//  - 若需要获取整个时间段的信息，请结合使用 Start 方法
//
// 使用建议：
//  - 确保 Period 实例有效且非零值，以避免返回无效的时间
//  - 在处理时间相关的逻辑时，优先使用内置方法以确保准确性和一致性
func (p Period) End() time.Time {
    return p[1]
}

// Duration 返回时间段的持续时间。
//
// 该方法计算并返回 Period 结构体中结束时间与开始时间之间的差值。
// 如果开始时间晚于结束时间，两者在创建时会自动交换以确保正确的顺序。
//
// 关键行为说明：
//  - 调用此方法不会改变 Period 的内部状态
//  - 若需要获取更细粒度的时间单位，请使用其他相关方法如 Days, Hours 等
//
// 使用建议：
//  - 确保 Period 实例有效且非零值，以避免返回无效的时间
//  - 在处理时间相关的逻辑时，优先使用内置方法以确保准确性和一致性
func (p Period) Duration() time.Duration {
    return p[1].Sub(p[0])
}

// Days 返回时间段的持续天数。
//
// 该方法通过计算时间段的总小时数并转换为天数来返回结果。
// 如果时间段小于一天，结果将被截断为整数天数。
//
// 关键行为说明：
//  - 时间段的持续时间由 Duration 方法计算
//  - 结果为整数天数，小数部分会被截断
func (p Period) Days() int {
    return int(p.Duration().Hours() / 24)
}

// Hours 返回时间段的持续小时数。
//
// 该方法通过计算时间段的总秒数并转换为小时数来返回结果。
// 如果时间段小于一小时，结果将被截断为整数小时数。
//
// 关键行为说明：
//  - 时间段的持续时间由 Duration 方法计算
//  - 结果为整数小时数，小数部分会被截断
func (p Period) Hours() int {
    return int(p.Duration().Hours())
}

// Minutes 返回时间段的持续分钟数。
//
// 该方法通过计算时间段的总秒数并转换为分钟数来返回结果。
// 如果时间段小于一分钟，结果将被截断为整数分钟数。
//
// 关键行为说明：
//  - 时间段的持续时间由 Duration 方法计算
//  - 结果为整数分钟数，小数部分会被截断
func (p Period) Minutes() int {
    return int(p.Duration().Minutes())
}

// Seconds 返回时间段的持续秒数。
//
// 该方法通过计算时间段的总纳秒数并转换为秒数来返回结果。
// 如果时间段小于一秒，结果将被截断为整数秒数。
//
// 关键行为说明：
//  - 时间段的持续时间由 Duration 方法计算
//  - 结果为整数秒数，小数部分会被截断
func (p Period) Seconds() int {
    return int(p.Duration().Seconds())
}

// Milliseconds 返回时间段的持续时间（以毫秒为单位）。
//
// 该方法通过计算时间段的总纳秒数并转换为毫秒数来返回结果。
// 如果时间段小于一毫秒，结果将被截断为整数毫秒数。
func (p Period) Milliseconds() int {
    return int(p.Duration().Milliseconds())
}

// Microseconds 将周期转换为微秒数。
//
// 该方法返回周期时长对应的微秒数。如果周期为零或负值，结果将反映相同的符号特性。
// 特别地，当周期非常短时，微秒数可能不足以精确表示实际时间间隔。
//
// 关键行为说明：
//  - 周期为零时，返回值也为零
//  - 负周期将返回负的微秒数
func (p Period) Microseconds() int {
    return int(p.Duration().Microseconds())
}

// Nanoseconds 将周期对象转换为纳秒数。
//
// 该方法返回当前周期持续时间对应的纳秒值，适用于需要更高精度时间计算的场景。对于零值周期，返回结果同样为0。
func (p Period) Nanoseconds() int {
    return int(p.Duration().Nanoseconds())
}

// IsZero 检查周期是否为零值。
//
// 该方法通过检查周期的开始和结束时间点是否都为零值来判断整个周期是否有效。如果两个时间点均为零，则返回 true，表示这是一个零值周期，否则返回 false。
//
// 关键行为说明：
//  - 零值周期意味着未定义或无效的时间范围
func (p Period) IsZero() bool {
    return p[0].IsZero() && p[1].IsZero()
}

// IsInvalid 检查周期是否无效。
//
// 通过判断周期的起始时间和结束时间是否为零值来确定该周期是否有效。
// 如果任一时间为零，则认为该周期是无效的，返回 true；否则返回 false。
func (p Period) IsInvalid() bool {
    return p[0].IsZero() || p[1].IsZero()
}

// Before 检查给定时间是否在当前周期之后。
//
// 参数 t 为要比较的时间点。方法通过内部存储的结束时间与 t 进行比较。
// 如果 t 在结束时间之前，则返回 true，表示 t 在当前周期内。
// 反之则返回 false，表明 t 已超出当前周期范围。
//
// 关键行为说明：
//  - 当 t 等于结束时间时，同样视为不在周期内
//
// 使用建议：
// 用于判断事件发生时间是否属于特定周期内，便于时间序列分析和处理。
func (p Period) Before(t time.Time) bool {
    return p[1].Before(t)
}

// After 检查给定时间是否在当前时间段之后。
//
// 该方法接收一个 time.Time 类型的参数 t，用于比较是否在时间段 p 之后。
// 如果 t 在 p 之后，则返回 true；否则返回 false。
// 注意，p[0] 表示时间段的起始时间。
func (p Period) After(t time.Time) bool {
    return p[0].After(t)
}

// Between 判断给定时间是否在周期内。
//
// 该方法接受一个时间点 t 作为参数，检查 t 是否位于由 p[0] 和 p[1] 定义的时间区间内。
// 如果 t 在 p[0]（含）和 p[1]（含）之间，则返回 true；否则返回 false。
//
// 关键行为说明：
//  - 当 p[0] 等于 p[1] 时，仅当 t 精确等于 p[0] 才返回 true
//  - 方法不考虑周期的循环特性，即 p[0] 不会自动视为晚于 p[1]
//
// 使用建议：
// 检查时间有效性前确保 p[0] 不晚于 p[1] 避免逻辑错误。
func (p Period) Between(t time.Time) bool {
    return (p[0].Before(t) || p[0].Equal(t)) && (p[1].After(t) || p[1].Equal(t))
}

// BetweenOrEqual 检查当前周期是否与给定周期重叠或相等。
//
// 该方法通过比较两个周期的起始和结束时间点来判断是否存在重叠或完全相同的情况。
// 如果当前周期的任意一个时间点在给定周期内，或者两个周期的任一端点相等，则返回true。
// 参数 t 表示要比较的目标周期。
//
// 关键行为说明：
//  - 当前周期或目标周期中任何一者为空时，结果为false
//  - 支持周期端点严格等于的情况
func (p Period) BetweenOrEqual(t Period) bool {
    return p.Between(t[0]) || p.Between(t[1]) || p[0].Equal(t[0]) || p[1].Equal(t[1])
}

// Overlap 检查两个时间段是否存在重叠。
//
// 该方法通过调用 BetweenOrEqual 方法判断两个时间段是否相互包含或边界相等来确定是否有重叠。
// 当任意一个时间段的开始时间或结束时间落在另一个时间段内时，认为存在重叠。
//
// 关键行为说明：
//  - 时间段的边界点也被视为有效范围
//  - 两个完全相同的时间段被视为完全重叠
//
// 使用建议：
// 确保输入的时间段是有效的，即开始时间不大于结束时间。
func (p Period) Overlap(t Period) bool {
    return p.BetweenOrEqual(t) || t.BetweenOrEqual(p)
}
