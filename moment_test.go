package chrono_test

import (
    "fmt"
    "github.com/kercylan98/chrono"
    "testing"
    "time"
)

func TestDelta(t *testing.T) {
    var curr = time.Now().AddDate(0, 0, -2)
    var a = chrono.StartOf(curr, chrono.UnitSaturday)
    fmt.Println(a)
    fmt.Println(a.AddDate(0, 0, 7))
    fmt.Println(a.AddDate(0, 0, -7))
}

func TestNextMoment(t *testing.T) {
    tests := []struct {
        name     string
        now      time.Time
        hour     int
        min      int
        sec      int
        expected time.Time
    }{
        {
            name:     "Before target moment",
            now:      time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
            hour:     15,
            min:      0,
            sec:      0,
            expected: time.Date(2023, 10, 1, 15, 0, 0, 0, time.Local),
        },
        {
            name:     "At target moment",
            now:      time.Date(2023, 10, 1, 15, 0, 0, 0, time.Local),
            hour:     15,
            min:      0,
            sec:      0,
            expected: time.Date(2023, 10, 2, 15, 0, 0, 0, time.Local),
        },
        {
            name:     "After target moment",
            now:      time.Date(2023, 10, 1, 16, 0, 0, 0, time.Local),
            hour:     15,
            min:      0,
            sec:      0,
            expected: time.Date(2023, 10, 2, 15, 0, 0, 0, time.Local),
        },
        {
            name:     "Midnight to next day",
            now:      time.Date(2023, 10, 1, 23, 59, 59, 0, time.Local),
            hour:     0,
            min:      0,
            sec:      0,
            expected: time.Date(2023, 10, 2, 0, 0, 0, 0, time.Local),
        },
        {
            name:     "Leap year",
            now:      time.Date(2024, 2, 28, 23, 59, 59, 0, time.Local),
            hour:     0,
            min:      0,
            sec:      0,
            expected: time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := chrono.NextMoment(tt.now, tt.hour, tt.min, tt.sec)
            if !result.Equal(tt.expected) {
                t.Errorf("NextMoment() = %v, want %v", result, tt.expected)
            }
        })
    }
}

func TestElapsed(t *testing.T) {
    tests := []struct {
        name     string
        now      time.Time
        hour     int
        min      int
        sec      int
        expected bool
    }{
        {
            name:     "Before target moment",
            now:      time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
            hour:     15,
            min:      0,
            sec:      0,
            expected: false,
        },
        {
            name:     "At target moment",
            now:      time.Date(2023, 10, 1, 15, 0, 0, 0, time.Local),
            hour:     15,
            min:      0,
            sec:      0,
            expected: false,
        },
        {
            name:     "After target moment",
            now:      time.Date(2023, 10, 1, 16, 0, 0, 0, time.Local),
            hour:     15,
            min:      0,
            sec:      0,
            expected: true,
        },
        {
            name:     "Midnight to next day",
            now:      time.Date(2023, 10, 1, 23, 59, 59, 0, time.Local),
            hour:     0,
            min:      0,
            sec:      0,
            expected: true,
        },
        {
            name:     "Leap year",
            now:      time.Date(2024, 2, 28, 23, 59, 59, 0, time.Local),
            hour:     0,
            min:      0,
            sec:      0,
            expected: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := chrono.Elapsed(tt.now, tt.hour, tt.min, tt.sec)
            if result != tt.expected {
                t.Errorf("Elapsed() = %v, want %v", result, tt.expected)
            }
        })
    }
}

func TestStartOf(t *testing.T) {
    tests := []struct {
        name     string
        now      time.Time
        unit     chrono.Unit
        expected time.Time
    }{
        {
            name:     "Nanosecond",
            now:      time.Date(2023, 10, 1, 12, 0, 0, 500, time.Local),
            unit:     chrono.UnitNanosecond,
            expected: time.Date(2023, 10, 1, 12, 0, 0, 500, time.Local),
        },
        {
            name:     "Microsecond",
            now:      time.Date(2023, 10, 1, 12, 0, 0, 123456789, time.Local),
            unit:     chrono.UnitMicrosecond,
            expected: time.Date(2023, 10, 1, 12, 0, 0, 123456000, time.Local),
        },
        {
            name:     "Millisecond",
            now:      time.Date(2023, 10, 1, 12, 0, 0, 123456789, time.Local),
            unit:     chrono.UnitMillisecond,
            expected: time.Date(2023, 10, 1, 12, 0, 0, 123000000, time.Local),
        },
        {
            name:     "Second",
            now:      time.Date(2023, 10, 1, 12, 0, 1, 123456789, time.Local),
            unit:     chrono.UnitSecond,
            expected: time.Date(2023, 10, 1, 12, 0, 1, 0, time.Local),
        },
        {
            name:     "Minute",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitMinute,
            expected: time.Date(2023, 10, 1, 12, 1, 0, 0, time.Local),
        },
        {
            name:     "Hour",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitHour,
            expected: time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
        },
        {
            name:     "Day",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitDay,
            expected: time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
        },
        {
            name:     "Week (Monday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitMonday,
            expected: time.Date(2023, 9, 25, 0, 0, 0, 0, time.Local),
        },
        {
            name:     "Week (Tuesday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitTuesday,
            expected: time.Date(2023, 9, 26, 0, 0, 0, 0, time.Local),
        },
        {
            name:     "Week (Wednesday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitWednesday,
            expected: time.Date(2023, 9, 27, 0, 0, 0, 0, time.Local),
        },
        {
            name:     "Week (Thursday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitThursday,
            expected: time.Date(2023, 9, 28, 0, 0, 0, 0, time.Local),
        },
        {
            name:     "Week (Friday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitFriday,
            expected: time.Date(2023, 9, 29, 0, 0, 0, 0, time.Local),
        },
        {
            name:     "Week (Saturday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitSaturday,
            expected: time.Date(2023, 9, 30, 0, 0, 0, 0, time.Local),
        },
        {
            name:     "Week (Sunday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitSunday,
            expected: time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
        },
        {
            name:     "Month",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitMonth,
            expected: time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
        },
        {
            name:     "Year",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 123456789, time.Local),
            unit:     chrono.UnitYear,
            expected: time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := chrono.StartOf(tt.now, tt.unit)
            if !result.Equal(tt.expected) {
                t.Errorf("StartOf() = %v, want %v", result, tt.expected)
            }
        })
    }
}

func TestFuture(t *testing.T) {
    tests := []struct {
        name     string
        now      time.Time
        hour     int
        min      int
        sec      int
        expected bool
    }{
        {
            name:     "Before target moment",
            now:      time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
            hour:     15,
            min:      0,
            sec:      0,
            expected: true,
        },
        {
            name:     "At target moment",
            now:      time.Date(2023, 10, 1, 15, 0, 0, 0, time.Local),
            hour:     15,
            min:      0,
            sec:      0,
            expected: true,
        },
        {
            name:     "After target moment",
            now:      time.Date(2023, 10, 1, 16, 0, 0, 0, time.Local),
            hour:     15,
            min:      0,
            sec:      0,
            expected: false,
        },
        {
            name:     "Midnight to next day",
            now:      time.Date(2023, 10, 1, 23, 59, 59, 0, time.Local),
            hour:     0,
            min:      0,
            sec:      0,
            expected: false,
        },
        {
            name:     "Leap year",
            now:      time.Date(2024, 2, 28, 23, 59, 59, 0, time.Local),
            hour:     0,
            min:      0,
            sec:      0,
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := chrono.Future(tt.now, tt.hour, tt.min, tt.sec)
            if result != tt.expected {
                t.Errorf("Future() = %v, want %v", result, tt.expected)
            }
        })
    }
}

func TestCeilDeltaDays(t *testing.T) {
    tests := []struct {
        name     string
        now      time.Time
        unit     chrono.Unit
        expected time.Time
    }{
        {
            name:     "Nanosecond",
            now:      time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
            unit:     chrono.UnitNanosecond,
            expected: time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
        },
        {
            name:     "Microsecond",
            now:      time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
            unit:     chrono.UnitMicrosecond,
            expected: time.Date(2023, 10, 1, 12, 0, 0, 999, time.Local),
        },
        {
            name:     "Millisecond",
            now:      time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
            unit:     chrono.UnitMillisecond,
            expected: time.Date(2023, 10, 1, 12, 0, 0, 999999, time.Local),
        },
        {
            name:     "Second",
            now:      time.Date(2023, 10, 1, 12, 0, 1, 0, time.Local),
            unit:     chrono.UnitSecond,
            expected: time.Date(2023, 10, 1, 12, 0, 1, 999999999, time.Local),
        },
        {
            name:     "Minute",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitMinute,
            expected: time.Date(2023, 10, 1, 12, 1, 59, 999999999, time.Local),
        },
        {
            name:     "Hour",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitHour,
            expected: time.Date(2023, 10, 1, 12, 59, 59, 999999999, time.Local),
        },
        {
            name:     "Day",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitDay,
            expected: time.Date(2023, 10, 1, 23, 59, 59, 999999999, time.Local),
        },
        {
            name:     "Week (Monday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitMonday,
            expected: time.Date(2023, 9, 25, 23, 59, 59, 999999999, time.Local),
        },
        {
            name:     "Week (Tuesday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitTuesday,
            expected: time.Date(2023, 9, 26, 23, 59, 59, 999999999, time.Local),
        },
        {
            name:     "Week (Wednesday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitWednesday,
            expected: time.Date(2023, 9, 27, 23, 59, 59, 999999999, time.Local),
        },
        {
            name:     "Week (Thursday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitThursday,
            expected: time.Date(2023, 9, 28, 23, 59, 59, 999999999, time.Local),
        },
        {
            name:     "Week (Friday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitFriday,
            expected: time.Date(2023, 9, 29, 23, 59, 59, 999999999, time.Local),
        },
        {
            name:     "Week (Saturday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitSaturday,
            expected: time.Date(2023, 9, 30, 23, 59, 59, 999999999, time.Local),
        },
        {
            name:     "Week (Sunday)",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitSunday,
            expected: time.Date(2023, 10, 1, 23, 59, 59, 999999999, time.Local),
        },
        {
            name:     "Month",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitMonth,
            expected: time.Date(2023, 10, 31, 23, 59, 59, 999999999, time.Local),
        },
        {
            name:     "Year",
            now:      time.Date(2023, 10, 1, 12, 1, 1, 0, time.Local),
            unit:     chrono.UnitYear,
            expected: time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.Local),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := chrono.EndOf(tt.now, tt.unit)
            if !result.Equal(tt.expected) {
                t.Errorf("StartOf() = %v, want %v", result, tt.expected)
            }
        })
    }
}
