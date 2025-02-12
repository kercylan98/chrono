package chrono

import "time"

// Truncate 将 x 截断为 m 的倍数
func Truncate(x, m int64) int64 {
	if m <= 0 {
		return x
	}
	return x - x%m
}

// ToMillisecond 将时间转换为毫秒
func ToMillisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
