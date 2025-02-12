package timing

import "time"

const (
	Nanosecond  = time.Nanosecond  // 纳秒
	Microsecond = time.Microsecond // 微秒
	Millisecond = time.Millisecond // 毫秒
	Second      = time.Second      // 秒
	Minute      = time.Minute      // 分钟
	Hour        = time.Hour        // 小时
	Day         = Hour * 24        // 天
	Week        = Day * 7          // 周
)
