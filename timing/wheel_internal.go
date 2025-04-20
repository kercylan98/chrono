package timing

import (
    "github.com/kercylan98/chrono"
    "github.com/kercylan98/chrono/timing/internal/delayqueue"
    "sync"
    "sync/atomic"
    "time"
)

var (
    _ Wheel = (*wheel)(nil)
)

func newWheelInternal(tw Wheel, config OptionsFetcher) wheelInternal {
    return &wheelInternalImpl{
        Wheel:  tw,
        config: config,
    }
}

type wheelInternal interface {
    // init 初始化时间轮
    init(startMs int64, queue *delayqueue.DelayQueue[bucket])

    // getConfig 获取时间轮的配置
    getConfig() OptionsFetcher

    // add 添加一个计时器
    add(timer Timer) bool

    // advanceClock 推进时间轮的时间
    advanceClock(expiration int64)

    // contract 履行任务
    contract(timer Timer)

    // refreshDelayQueue 刷新延迟队列，避免长时间无效挂起
    refreshDelayQueue()
}

type wheelInternalImpl struct {
    Wheel
    config       OptionsFetcher                 // 时间轮的配置
    overflow     Wheel                          // 溢出轮
    overflowLock sync.RWMutex                   // 溢出轮锁
    buckets      []bucket                       // 时间轮的桶
    queue        *delayqueue.DelayQueue[bucket] // 延迟队列
    current      int64                          // 毫秒级当前时间
    interval     int64                          // 时间轮的间隔时间
}

func (t *wheelInternalImpl) init(startMs int64, queue *delayqueue.DelayQueue[bucket]) {
    if startMs == 0 {
        startMs = chrono.ToMillisecond(time.Now())
    }
    tick := t.getConfig().FetchTick()
    size := t.getConfig().FetchSize()

    t.current = chrono.Truncate(startMs, tick)
    t.interval = tick * size
    t.buckets = make([]bucket, size)

    if queue == nil {
        queue = delayqueue.New(int(size), func() int64 {
            return chrono.ToMillisecond(time.Now())
        }, func(bucket bucket) {
            t.advanceClock(bucket.getExpiration())
            bucket.flush(t.contract)
        })
    }
    t.queue = queue

    for i := range t.buckets {
        t.buckets[i] = newBucket(t)
    }
}

func (t *wheelInternalImpl) getConfig() OptionsFetcher {
    return t.config
}

func (t *wheelInternalImpl) contract(timer Timer) {
    if timer.Stopped() {
        return
    }
    if !t.add(timer) {
        // 计时器已经过期，直接执行
        go t.getConfig().FetchExecutor().Execute(timer.getTask())
    }
}

func (t *wheelInternalImpl) add(timer Timer) bool {
    // 获取时间轮当前时间和下一个刻度时间，以及待添加的计时器的到期时间
    current := atomic.LoadInt64(&t.current)
    tick := t.getConfig().FetchTick()
    expiration := timer.getExpiration()
    if expiration < current+tick {
        // 计时器已经过期
        return false
    } else if expiration < current+t.interval {
        // 计算计时器位于时间轮的哪个刻度，然后获取对应的桶
        b := t.buckets[expiration/tick%t.getConfig().FetchSize()]
        b.add(timer)
        if b.setExpiration(expiration) {
            // 如果桶的过期时间发生变化，需要重新调度桶
            t.queue.Add(b, b.getExpiration())
        }
        return true
    } else {
        // 超出区间。将其放入溢流轮中
        t.overflowLock.Lock()
        defer t.overflowLock.Unlock()
        if t.overflow == nil {
            config := NewConfig().
                withTick(t.interval).
                WithSize(int(t.getConfig().FetchSize())).
                WithExecutor(t.getConfig().FetchExecutor())
            t.overflow = GetBuilder().build(current, t.queue, config)
        }
        return t.overflow.add(timer)
    }
}

func (t *wheelInternalImpl) advanceClock(expiration int64) {
    currentTime := atomic.LoadInt64(&t.current)
    tick := t.getConfig().FetchTick()
    if expiration >= currentTime+tick {
        // 当给定的时间超出当前时间轮的间隔时推进时间轮的时间
        currentTime = chrono.Truncate(expiration, tick)
        atomic.StoreInt64(&t.current, currentTime)

        // 如果溢出时间轮存在，则同时推进溢出时间轮的时间
        t.overflowLock.RLock()
        defer t.overflowLock.RUnlock()
        if t.overflow != nil {
            t.overflow.advanceClock(currentTime)
        }
    }
}

func (t *wheelInternalImpl) refreshDelayQueue() {
    t.queue.Refresh()
}
