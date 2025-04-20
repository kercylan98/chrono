package chrono

import "time"

// ToMillisecond 将时间对象转换为毫秒时间戳。
//
// t 参数为需要转换的时间对象，函数将该时间对象的纳秒数除以 1000000 得到毫秒时间戳。
//
// 关键行为说明：
//  - 转换结果为 int64 类型的毫秒时间戳
//  - 时间精度依赖于输入时间对象的精度
func ToMillisecond(t time.Time) int64 {
    return t.UnixNano() / int64(time.Millisecond)
}

// ToTime 将给定的毫秒数转换为UTC时间。
//
// mill 参数表示自 Unix 纪元以来的毫秒数，函数将此值转换为对应的时间对象。
//
// 关键行为说明：
//  - 输入为0时返回 Unix 纪元开始时刻
//  - 负值输入同样有效，表示纪元前的时间
//
// 使用建议：
//  - 确保输入值在 int64 范围内
//  - 注意处理可能的时区差异
func ToTime(mill int64) time.Time {
    return time.Unix(0, mill*int64(time.Millisecond)).UTC()
}

// Truncate 将 x 以 m 为单位进行截断，返回最接近 x 且不大于 x 的 m 的倍数。
//
// 参数 x 表示要截断的整数值，m 表示截断的模数。当 m 小于等于 0 时，函数直接返回 x。
// 函数通过计算 x - x % m 来实现截断操作，确保结果是 m 的倍数且不超过 x。
//
// 关键行为说明：
//  - 当 m <= 0 时，直接返回 x 不做任何修改
//  - 截断操作基于数学模运算，适用于需要对齐到特定间隔的场景
func Truncate(x, m int64) int64 {
    if m <= 0 {
        return x
    }
    return x - x%m
}
