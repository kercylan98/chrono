package timing_test

import (
	"context"
	"github.com/kercylan98/chrono/timing"
	"testing"
	"time"
)

func TestWheel_Loop(t *testing.T) {
	var tests = []struct {
		name     string
		interval time.Duration
		times    int
	}{
		{"100ms-10", 100 * time.Millisecond, 10},
		{"500ms-10", 500 * time.Millisecond, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := timing.New()
			startAt := time.Now()
			finish := make([]time.Time, 0, tt.times)
			wait := make(chan time.Time, 1)
			ctx, cancel := context.WithTimeout(context.Background(), tt.interval*time.Duration(tt.times+1)+time.Second)
			defer cancel()
			tw.Loop(0, timing.NewLoopTask(tt.interval, tt.times, timing.TaskFn(func() {
				wait <- time.Now()
			})))
			for len(finish) != tt.times {
				select {
				case <-ctx.Done():
					t.Fatalf("timeout")
				case finishTime := <-wait:
					finish = append(finish, finishTime)
					if len(finish) == tt.times {
						close(wait)
					}
				}
			}

			if len(finish) != tt.times {
				t.Fatalf("loop: want %d, got %d", tt.times, len(finish))
			}

			finish = append([]time.Time{startAt.Add(-tt.interval)}, finish...)

			// 验证时间
			for i := 1; i < len(finish); i++ {
				before := finish[i-1]
				curr := finish[i]
				delta := curr.Sub(before) - tt.interval
				if delta < -5*time.Millisecond || delta > 5*time.Millisecond {
					t.Errorf("FAIL loop[%d]: want +-5, got %s", i, delta)
				} else {
					t.Logf("PASS loop[%d]: pass +-5, got %s", i, delta)
				}

			}
		})
	}
}

func TestWheel(t *testing.T) {
	tw := timing.New()

	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
	for _, d := range durations {
		t.Run("", func(t *testing.T) {
			exitC := make(chan time.Time)

			start := time.Now().UTC()
			tw.After(d, timing.TaskFn(func() {
				exitC <- time.Now().UTC()
			}))

			got := (<-exitC).Truncate(time.Millisecond)
			min := start.Add(d).Truncate(time.Millisecond)

			err := 5 * time.Millisecond
			if got.Before(min) || got.After(min.Add(err)) {
				t.Errorf("Timer(%s) expiration: want [%s, %s], got %s", d, min, min.Add(err), got)
			}
		})
	}
}
