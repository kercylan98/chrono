package delayqueue

func newPriorityQueueItem[T any](value T, priority int64) *priorityQueueItem[T] {
	return &priorityQueueItem[T]{
		Value:    value,
		Priority: priority,
	}
}

type priorityQueueItem[T any] struct {
	Value    T
	Priority int64
}
