package chrono

import "time"

// Unit 定义了时间单位，用于表示时间间隔或持续时间。
//
// 该类型通常与时间相关的操作一起使用，例如定时任务的调度、延迟执行等。支持的时间单位包括秒、毫秒等。
type Unit int

const (
    UnitSunday      Unit = 0                 // UnitSunday 表示星期天，用于定义以星期为基准的时间间隔或持续时间。
    UnitMonday      Unit = 10                // UnitMonday 表示星期一，用于定义以星期为基准的时间间隔或持续时间。
    UnitTuesday     Unit = 20                // UnitTuesday 表示星期二，用于定义以星期为基准的时间间隔或持续时间。
    UnitWednesday   Unit = 30                // UnitWednesday 表示星期三，用于定义以星期为基准的时间间隔或持续时间。
    UnitThursday    Unit = 40                // UnitThursday 表示星期四，用于定义以星期为基准的时间间隔或持续时间。
    UnitFriday      Unit = 50                // UnitFriday 表示星期五，用于定义以星期为基准的时间间隔或持续时间。
    UnitSaturday    Unit = 60                // UnitSaturday 表示星期六，用于定义以星期为基准的时间间隔或持续时间。
    UnitNanosecond       = Unit(Nanosecond)  // UnitNanosecond 定义了纳秒时间单位，适用于需要高精度时间控制的场景。
    UnitMicrosecond      = Unit(Microsecond) // UnitMicrosecond 定义了微秒时间单位，适用于需要较高精度时间控制的场景。
    UnitMillisecond      = Unit(Millisecond) // UnitMillisecond 定义了毫秒时间单位，适用于需要中等精度时间控制的场景。
    UnitSecond           = Unit(Second)      // UnitSecond 定义了秒时间单位，适用于需要秒级精度时间控制的场景。
    UnitMinute           = Unit(Minute)      // UnitMinute 定义了分钟时间单位，适用于需要以分钟为精度的时间控制场景。
    UnitHour             = Unit(Hour)        // UnitHour 定义了小时时间单位，适用于需要以小时为精度的时间控制场景。
    UnitDay              = Unit(Day)         // UnitDay 定义了天时间单位，适用于需要以天为精度的时间控制场景。
    UnitWeek             = Unit(Week)        // UnitWeek 定义了周时间单位，适用于需要以周为精度的时间控制场景。
    UnitMonth            = Unit(Week * 30)   // UnitMonth 表示月时间单位，用于定义以月份为基准的时间间隔或持续时间。
    UnitYear             = UnitMonth * 12    // UnitYear 表示年时间单位，用于定义长时间间隔或持续时间。

)

const (
    // Nanosecond 表示时间单位纳秒，用于时间测量和计算。
    Nanosecond = time.Nanosecond

    // Microsecond 表示时间单位微秒，用于精确的时间测量和计算。
    //
    // 参数说明：
    //  - 该常量定义了时间的微秒单位，适用于需要高精度时间处理的场景。
    //  - 在进行时间转换或比较时，可以使用此常量确保时间单位的一致性。
    //
    // 关键行为说明：
    //  - 使用时需注意系统时钟的精度限制，实际操作中可能存在细微偏差。
    Microsecond = time.Microsecond

    // Millisecond 表示时间单位毫秒，用于时间测量和计算。
    //
    // 该常量定义了时间的毫秒单位，适用于需要较高精度时间处理的场景。
    // 在进行时间转换或比较时，可以使用此常量确保时间单位的一致性。
    //
    // 关键行为说明：
    //  - 使用时需注意系统时钟的精度限制，实际操作中可能存在细微偏差。
    Millisecond = time.Millisecond

    // Second 表示一秒的时间持续。
    //
    // 该常量用于表示时间间隔，等同于 time.Second。在需要指定以秒为单位的时间时使用。
    Second = time.Second

    // Minute 表示一分钟的时间持续，等同于 60 秒。
    //
    // 该常量用于时间相关的计算与操作中，作为时间单位使用。
    Minute = time.Minute

    // Hour 表示一个小时的时间持续时间。
    //
    // 该常量用于表示时间间隔，其值等同于 time.Hour。在定义较长的时间周期时，
    // 可以通过此常量进行计算，例如 Day = Hour * 24。
    Hour = time.Hour

    // Day 表示一天的时间持续时间。
    //
    // 该常量基于 Hour 定义，表示 24 小时的时间间隔。适用于需要以天为单位的时间计算场景。
    Day = Hour * 24

    // Week 表示一周的时间持续时间。
    //
    // 该常量基于 Day 定义，表示 7 天的时间间隔。适用于需要以周为单位的时间计算场景。
    Week = Day * 7
)

// zero 表示时间的零值，通常用于初始化或比较。
//
// 该变量定义了一个没有任何有效时间信息的时间点，可用于判断其他时间是否被明确设置。在时间相关的逻辑中，用作默认值或哨兵值以简化代码实现。
//
// 关键行为说明：
//  - 与任何非零时间比较时总是更早
//  - 通过 time.IsZero 方法可以方便地检测一个时间是否为零值
var zero = time.Time{}
