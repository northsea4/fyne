//go:build !wasm && !test_web_driver

package glfw

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-gl/glfw/v3.4/glfw"
)

var nativeWakeSerial atomic.Uint64
var nativeLoopWaiting atomic.Bool
var nativeWakeMu sync.Mutex
var nativeWakeReady bool // protected with nativeWakeMu through PostEmptyEvent

// WakeUp signals new application work without keeping an idle timer alive.
func (d *gLDriver) WakeUp() { wakeUpDriver() }
func wakeUpDriver() {
	nativeWakeSerial.Add(1)
	if !nativeLoopWaiting.Load() {
		return
	}
	nativeWakeMu.Lock()
	if nativeWakeReady {
		glfw.PostEmptyEvent()
	}
	nativeWakeMu.Unlock()
}
func stopNativeWake() { nativeWakeMu.Lock(); nativeWakeReady = false; nativeWakeMu.Unlock() }

// UsesNativeEventWait lets platform hosts avoid adding their own idle polling.
func (d *gLDriver) UsesNativeEventWait() bool { return true }

func (d *gLDriver) hasVisibleWork() bool {
	for _, win := range d.windowList() {
		w := win.(*window)
		if w.closing || !w.visible {
			continue
		}
		if w.shouldExpand || w.canvas.IsDirty() {
			return true
		}
	}
	return false
}
func (d *gLDriver) maintainWindows() {
	for i := 0; i < len(d.windows); i++ {
		w := d.windows[i].(*window)
		w.ensurePositionProcessed()
		if w.viewport == nil {
			continue
		}
		if w.viewport.ShouldClose() {
			d.destroyWindow(w, i)
			i--
			continue
		}
		if w.shouldExpand && !w.fullScreen {
			w.fitContent()
			expand := w.shouldExpand
			w.shouldExpand = false
			if expand {
				w.viewport.SetSize(w.shouldWidth, w.shouldHeight)
			}
		}
	}
}

// Input callbacks stay synchronous. Work functions execute as soon as woken;
// only animation/painting are paced. Continuous input never moves the deadline.
func (d *gLDriver) runIdleLoop() bool {
	nativeWakeMu.Lock()
	nativeWakeReady = true
	nativeWakeMu.Unlock()
	nextFrame := time.Now()
	const frameInterval = time.Second / 60
	for running.Load() {
		observed := nativeWakeSerial.Load()
		// A snapshot prevents continuous producers from starving native input/frame.
		runWork(funcQueue.take(), true)
		if !running.Load() {
			break
		}
		d.pollEvents()
		d.maintainWindows()
		if !running.Load() {
			break
		}
		active := d.animation.HasAnimations() || d.hasVisibleWork()
		now := time.Now()
		if active && !now.Before(nextFrame) {
			d.animation.TickAnimations()
			d.drawSingleFrame()
			d.maintainWindows()
			nextFrame = now.Add(frameInterval)
		}
		if !running.Load() {
			break
		}
		if funcQueue.pending() {
			continue
		}
		active = d.animation.HasAnimations() || d.hasVisibleWork()
		nativeLoopWaiting.Store(true)
		// Close the gap between inspecting queues/state and entering native wait.
		// A producer either changes this serial or posts while waiting is true.
		if nativeWakeSerial.Load() != observed {
			nativeLoopWaiting.Store(false)
			continue
		}
		if active {
			delay := time.Until(nextFrame)
			if delay > 0 {
				glfw.WaitEventsTimeout(delay.Seconds())
			}
		} else {
			glfw.WaitEvents()
		}
		nativeLoopWaiting.Store(false)
	}
	nativeLoopWaiting.Store(false)
	return true
}
