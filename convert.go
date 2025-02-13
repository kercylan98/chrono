package chrono

import "time"

// TimeToMillisecond 将时间转换为 int64 类型的毫秒
func TimeToMillisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// MillisecondToTime 将 int64 类型的毫秒转换为时间
func MillisecondToTime(mill int64) time.Time {
	return time.Unix(0, mill*int64(time.Millisecond)).UTC()
}

// Truncate 将 x 截断为 m 的倍数
//   - 如果 m 小于等于 0，则返回 x
func Truncate(x, m int64) int64 {
	if m <= 0 {
		return x
	}
	return x - x%m
}
