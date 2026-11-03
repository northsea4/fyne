//go:build !no_glfw && !mobile

package glfw

import (
	"sync"
	"testing"
	"time"
)

func TestMainWorkQueueCloseReleasesConcurrentWaiters(t *testing.T) {
	q := newMainWorkQueue()
	const producers = 128
	var wg sync.WaitGroup
	ready := make(chan struct{})
	for i := 0; i < producers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			done := make(chan struct{}, 1)
			if q.push(funcData{f: func() { t.Error("cancelled work executed") }, done: done}) {
				<-done
			}
		}()
	}
	close(ready)
	runWork(q.close(), false)
	finished := make(chan struct{})
	go func() { wg.Wait(); close(finished) }()
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("waiter stranded by concurrent close")
	}
	if q.push(funcData{}) {
		t.Fatal("accepted work after close")
	}
}

func TestIdleLoopWakesForBackgroundWork(t *testing.T) {
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		done := make(chan struct{})
		go runOnMain(func() { close(done) })
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatalf("idle wake lost at request %d", i)
		}
	}
}
