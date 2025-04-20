package chrono

import (
    "time"
)

// NextMoment 计算并返回指定时间点在今天或明天的时刻。
//
// now 参数表示当前时间，用于与目标时刻进行比较。hour, min, sec 参数共同定义了具体的目标时刻。
// 如果目标时刻已经过去（即now大于等于moment），则自动调整为目标时刻的次日同一时间。
//
// 关键行为说明：
//  - 当前时间晚于或等于目标时刻时，返回值为次日同一时刻
//  - 输入的时间参数不受夏令时影响，始终基于本地时区计算
//
// 使用建议：
//  - 确保输入的时间参数合理，避免出现无效时间组合
func NextMoment(now time.Time, hour, min, sec int) time.Time {
    moment := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, time.Local)
    // 如果要检查的时刻已经过了，则返回明天的这个时刻
    if now.After(moment) || now.Equal(moment) {
        moment = moment.AddDate(0, 0, 1)
    }
    return moment
}

// Elapsed 判断给定的时刻是否已经过去。
//
// 参数 now 表示当前时间，hour、min 和 sec 分别表示指定时刻的小时、分钟和秒。
// 如果指定的时刻已经过去，则返回 true；否则返回 false。
//
// 关键行为说明：
//  - 当前时间与指定时刻相同视为已过期
//  - 指定时刻基于当前日期计算，不考虑跨天情况
//
// 使用建议：
//  - 用于判断特定时间点是否已经到达或超过
//  - 不适用于需要跨天处理的时间比较
func Elapsed(now time.Time, hour, min, sec int) bool {
    moment := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
    return now.After(moment)
}

// Future 判断给定的时刻是否尚未到达。
//
// 参数 now 表示当前时间，hour、min 和 sec 分别表示指定时刻的小时、分钟和秒。
// 如果指定的时刻还未到达，则返回 true；否则返回 false。
//
// 关键行为说明：
//  - 当前时间与指定时刻相同视为已到达
//  - 指定时刻基于当前日期计算，不考虑跨天情况
//
// 使用建议：
//  - 用于判断特定时间点是否还未到达
//  - 不适用于需要跨天处理的时间比较
func Future(now time.Time, hour, min, sec int) bool {
    return !Elapsed(now, hour, min, sec)
}

// StartOf 根据给定的时间单位，计算并返回时间 t 的起始点。
//
// 参数 t 为需要计算的时间点。unit 用于指定时间的度量单位，如小时、天等。
// 当 unit 为零或负值时，默认使用一天作为时间单位。
// 返回的时间是根据指定单位对 t 进行向下取整后的结果。
//
// 关键行为说明：
//  - 如果 t 本身已经是单位的起点，则直接返回 t
//  - 对于定义外的单位，函数会抛出异常
//
// 使用建议：
// 确保传递给 unit 的是一个标准的时间单位，例如 UnitDay、 UnitHour 等。
// 避免使用自定义的时间间隔以防止潜在的错误
func StartOf(t time.Time, unit Unit) time.Time {
    if unit <= 0 {
        unit = UnitDay
    }
    switch unit {
    case UnitNanosecond:
        return t.Truncate(Nanosecond)
    case UnitMicrosecond:
        return t.Truncate(Microsecond)
    case UnitMillisecond:
        return t.Truncate(Millisecond)
    case UnitSecond:
        return t.Truncate(Second)
    case UnitMinute:
        return t.Truncate(Minute)
    case UnitHour:
        return t.Truncate(Hour)
    case UnitDay:
        return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
    case UnitWeek, UnitMonday, UnitTuesday, UnitWednesday, UnitThursday, UnitFriday, UnitSaturday, UnitSunday:
        unit /= 10
        t = StartOf(t, UnitDay)
        tw := t.Weekday()
        if tw == 0 {
            tw = 7
        }
        d := 1 - int(tw)
        switch unit {
        case UnitSunday:
            d += 6
        default:
            if unit == UnitWeek {
                unit = UnitMonday
            }
            d += int(unit) - 1
        }
        return t.AddDate(0, 0, d)
    case UnitMonth:
        return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
    case UnitYear:
        return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
    default:
        panic("unsupported time unit")
    }
}

// EndOf 根据给定的时间单位，计算并返回时间 t 的结束点。
//
// 参数 t 为需要计算的时间点。unit 用于指定时间的度量单位，如小时、天等。
// 当 unit 为零或负值时，默认使用一天作为时间单位。
// 返回的时间是根据指定单位对 t 进行向上取整后的结果。
//
// 关键行为说明：
//  - 如果 t 本身已经是单位的终点，则直接返回 t
//  - 对于定义外的单位，函数会抛出异常
//
// 使用建议：
// 确保传递给 unit 的是一个标准的时间单位，例如 UnitDay、 UnitHour 等。
// 避免使用自定义的时间间隔以防止潜在的错误
func EndOf(t time.Time, unit Unit) time.Time {
    if unit <= 0 {
        unit = UnitDay
    }
    switch unit {
    case UnitNanosecond:
        return t.Truncate(Nanosecond)
    case UnitMicrosecond:
        return t.Truncate(Microsecond).Add(Microsecond - 1)
    case UnitMillisecond:
        return t.Truncate(Millisecond).Add(Millisecond - 1)
    case UnitSecond:
        return t.Truncate(Second).Add(Second - 1)
    case UnitMinute:
        return t.Truncate(Minute).Add(Minute - 1)
    case UnitHour:
        return t.Truncate(Hour).Add(Hour - 1)
    case UnitDay:
        return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
    case UnitWeek, UnitMonday, UnitTuesday, UnitWednesday, UnitThursday, UnitFriday, UnitSaturday, UnitSunday:
        unit /= 10
        t = EndOf(t, UnitDay)
        tw := t.Weekday()
        if tw == 0 {
            tw = 7
        }
        d := 1 - int(tw)
        switch unit {
        case UnitSunday:
            d += 6
        default:
            if unit == UnitWeek {
                unit = UnitMonday
            }
            d += int(unit) - 1
        }
        return EndOf(t.AddDate(0, 0, d), UnitDay)
    case UnitMonth:
        return StartOf(t, unit).AddDate(0, 1, 0).Add(-time.Nanosecond)
    case UnitYear:
        return StartOf(t, unit).AddDate(1, 0, 0).Add(-time.Nanosecond)
    default:
        panic("unsupported time unit")
    }
}

// Zero 返回表示时间零值的Time对象，用于初始化或比较。
func Zero() time.Time {
    return zero
}

// Max 返回两个时间点中较晚的那个。
//
// 该函数接受两个 time.Time 类型参数，比较它们的时间先后，并返回较晚的一个。如果两个时间相等，则返回任一参数。
//
// 参数说明：
// - t1: 第一个待比较的时间点
// - t2: 第二个待比较的时间点
//
// 关键行为说明：
//  - 如果 t1 和 t2 相等，函数将返回 t1
func Max(t1, t2 time.Time) time.Time {
    if t1.After(t2) {
        return t1
    }
    return t2
}

// Min 返回两个时间点中较早的一个。
//
// t1 和 t2 是需要比较的两个时间点。如果 t1 在 t2 之前，则返回 t1，否则返回 t2。
//
// 关键行为说明：
//  - 如果两个时间相等，将返回第一个参数 t1
//  - 时间点的比较基于 Go 的 time.Time 类型定义
func Min(t1, t2 time.Time) time.Time {
    if t1.Before(t2) {
        return t1
    }
    return t2
}

// SmallerFirst 返回两个时间中较早的一个作为第一个返回值。
//
// 该函数接收两个 time.Time 类型的参数 t1 和 t2，比较它们的时间先后顺序。
// 如果 t1 在 t2 之前，则返回 t1 作为第一个结果，t2 作为第二个结果；反之亦然。
// 这种机制确保了无论输入的时间顺序如何，总是能获得一个有序的时间对。
//
// 关键行为说明：
// - 当两个时间相等时，返回顺序与输入保持一致
//
// 使用建议：
// - 可用于需要按时间排序的场景，如日志记录或事件处理
func SmallerFirst(t1, t2 time.Time) (time.Time, time.Time) {
    if t1.Before(t2) {
        return t1, t2
    }
    return t2, t1
}

// SmallerLast 比较两个时间点，返回时序上靠后的和靠前的时间。
//
// 该函数接收两个 time.Time 类型参数 t1 和 t2，比较它们的先后顺序。
// 如果 t1 在 t2 之前，则返回 (t2, t1)，否则返回 (t1, t2)。
// 这种处理方式确保了第一个返回值总是两者中更晚的那个时间点。
//
// 参数说明：
// - t1: 第一个时间点
// - t2: 第二个时间点
//
// 关键行为说明：
//  - 函数直接使用 time.Time 的 Before 方法进行比较，确保了时间点判断的准确性
//  - 若两时间相同，将原样返回输入参数
//
// 使用建议：
//  - 用于需要确定时间序列关系的场景，如日志排序、事件触发顺序等
func SmallerLast(t1, t2 time.Time) (time.Time, time.Time) {
    if t1.Before(t2) {
        return t2, t1
    }
    return t1, t2
}

// Delta 计算两个时间点之间的时间差。
//
// 参数 t1 和 t2 分别表示要比较的两个时间点。函数会自动判断哪个时间点更早，并计算它们之间的时间差。如果 t1 在 t2 之前，返回值为 t2 减去 t1 的时间差；反之，则返回 t1 减去 t2 的时间差。这样确保了无论参数顺序如何，返回的时间差总是非负的。
//
// 关键行为说明：
//  - 时间差以 time.Duration 类型返回，单位为纳秒
//  - 当两个时间点相同时，返回的时间差为零
func Delta(t1, t2 time.Time) time.Duration {
    if t1.Before(t2) {
        return t2.Sub(t1)
    }
    return t1.Sub(t2)
}

// MonthDays 返回给定时间的月份天数。
//
// 参数 t 影响函数行为，它决定了返回哪个月份的天数。对于非二月，特定月份有固定的天数：4、6、9 和 11 月为 30 天，其他月份为 31 天。对于二月，根据年份是否为闰年来决定天数：普通年份 28 天，闰年 29 天。
//
// 关键行为说明：
//  - 函数仅依赖于提供的 time.Time 类型参数来确定月份和年份
//  - 闰年的判断基于格里高利历规则
//
// 使用建议：
// 确保输入的时间值是有效的，以避免意外的行为。
func MonthDays(t time.Time) int {
    year, month, _ := t.Date()
    if month != 2 {
        if month == 4 || month == 6 || month == 9 || month == 11 {
            return 30
        }
        return 31
    }
    if ((year%4 == 0) && (year%100 != 0)) || year%400 == 0 {
        return 29
    }
    return 28
}
