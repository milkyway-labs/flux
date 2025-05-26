package indexer

import (
	"context"
	"time"
)

// Queue is a generic, thread-safe, buffered queue for transferring
// objects between goroutines.
type Queue[T any] struct {
	channel chan T
}

// NewQueue creates and returns a new buffered queue with the specified size.
func NewQueue[T any](size uint32) *Queue[T] {
	return &Queue[T]{
		channel: make(chan T, size),
	}
}

// Enqueue inserts a new value into the queue.
// If the queue is full, this call will block until space becomes available.
func (q *Queue[T]) Enqueue(value T) {
	q.channel <- value
}

// EnqueueWithContext attempts to insert a new value into the queue.
// If the queue is full, this call will block until space becomes available or the context is canceled.
// Returns true if the value was enqueued, or false if the context was canceled first.
func (q *Queue[T]) EnqueueWithContext(ctx context.Context, value T) bool {
	select {
	case <-ctx.Done():
		return false
	case q.channel <- value:
		return true
	}
}

// DelayedEnqueue schedules the insertion of a value into the queue after the specified delay.
// If the context is canceled before the delay elapses, the enqueue operation is aborted.
func (q *Queue[T]) DelayedEnqueue(ctx context.Context, delay time.Duration, value T) {
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
			q.EnqueueWithContext(ctx, value)
		}
	}()
}

// Dequeue removes and returns a value from the queue.
// If the queue is empty, this call will block until a value is available or the queue is closed.
// Returns (value, true) if successful, or (zero, false) if the queue has been closed and emptied.
func (q *Queue[T]) Dequeue() (T, bool) {
	value, ok := <-q.channel
	return value, ok
}

// ContextDequeue attempts to dequeue a value from the queue only if the context has not been canceled.
// If the context is canceled, returns (zero, false) immediately.
// Otherwise, behaves like Dequeue.
func (q *Queue[T]) ContextDequeue(ctx context.Context) (T, bool) {
	var zero T
	select {
	case <-ctx.Done():
		return zero, false
	default:
		return q.Dequeue()
	}
}

// Close closes the queue, indicating that no more values will be enqueued.
// After closing, all future Dequeue operations will return (zero, false) once the queue is empty.
func (q *Queue[T]) Close() {
	close(q.channel)
}
