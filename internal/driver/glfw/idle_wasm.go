//go:build wasm || test_web_driver

package glfw

// Browsers retain their existing event/timer implementation in this fork.
func (d *gLDriver) runIdleLoop() bool { return false }
func (d *gLDriver) WakeUp()           {}
func wakeUpDriver()                   {}
func stopNativeWake()                 {}
