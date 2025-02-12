package timing_test

import (
	"github.com/kercylan98/chrono/timing"
	"testing"
	"time"
)

func genD(i int) time.Duration {
	return time.Duration(i%10000) * time.Millisecond
}

func BenchmarkWheel_AfterFunc(b *testing.B) {
	tw := timing.New()

	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1m", 100000},
		{"N-5m", 500000},
		{"N-10m", 1000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]timing.Timer, c.N)
			for i := 0; i < len(base); i++ {
				base[i] = tw.AfterFunc(genD(i), func() {})
			}
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				tw.AfterFunc(time.Second, func() {}).Stop()
			}

			b.StopTimer()
			for i := 0; i < len(base); i++ {
				base[i].Stop()
			}
		})
	}
}
