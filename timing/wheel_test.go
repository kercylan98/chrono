package timing_test

import (
    "fmt"
    "github.com/kercylan98/chrono/timing"
    "testing"
    "time"
)

func TestWheel_After(t *testing.T) {
    tw := timing.New()
    tw.Loop(0, timing.NewForeverLoopTask(-124, timing.TaskFN(func() {
        fmt.Println(1)
    })))

    time.Sleep(time.Second)
}
