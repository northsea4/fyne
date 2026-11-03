package glfw

import "sync"

// mainWorkQueue publishes work before the producer wakes the native loop.
// A channel forwarded by another goroutine cannot provide that ordering.
type mainWorkQueue struct {
	mu     sync.Mutex
	items  []funcData
	closed bool
	notify chan struct{}
}

func newMainWorkQueue() *mainWorkQueue { return &mainWorkQueue{notify: make(chan struct{}, 1)} }
func (q *mainWorkQueue) push(item funcData) bool {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return false
	}
	q.items = append(q.items, item)
	select {
	case q.notify <- struct{}{}:
	default:
	}
	q.mu.Unlock()
	return true
}
func (q *mainWorkQueue) take() []funcData {
	q.mu.Lock()
	defer q.mu.Unlock()
	items := q.items
	q.items = nil
	return items
}
func (q *mainWorkQueue) pending() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items) > 0
}

// Close and admission use the same lock. Every accepted waiter is returned to
// the caller for cancellation; no producer can publish behind the close drain.
func (q *mainWorkQueue) close() []funcData {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	items := q.items
	q.items = nil
	return items
}
func runWork(items []funcData, execute bool) {
	for _, item := range items {
		if execute {
			item.f()
		}
		if item.done != nil {
			item.done <- struct{}{}
		}
	}
}
