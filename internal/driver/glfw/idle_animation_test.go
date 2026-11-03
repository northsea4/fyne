//go:build !wasm && !test_web_driver && !no_glfw && !mobile

package glfw

import (
	"sync/atomic"
	"testing"
	"time"

	"fyne.io/fyne/v2"
)

func TestIdleLoopAnimationWakesAndStops(t *testing.T) {
	var ticks atomic.Int64
	done := make(chan struct{})
	a := fyne.NewAnimation(150*time.Millisecond, func(v float32) {
		ticks.Add(1)
		if v == 1 {
			close(done)
		}
	})
	d.StartAnimation(a)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("animation failed to wake idle driver")
	}
	n := ticks.Load()
	if n < 2 || n > 15 {
		t.Fatalf("animation ticks not paced near 60Hz: %d", n)
	}
	var active bool
	runOnMain(func() { active = d.animation.HasAnimations() })
	if active {
		t.Fatal("completed animation kept the loop active")
	}
}
