//go:build darwin && !cddrop_no_native_present

package glfw

import "fyne.io/fyne/v2"

var nativeFrameDrawing bool

// RenderWindowNow submits a dirty window while AppKit owns a nested tracking
// loop. The caller must already be on the Fyne/GLFW main thread. It does not poll
// input, drain queued work, run animations or recursively enter the driver loop.
// CDDrop extension to Fyne 2.8.2; see CDDROP_PATCH.md.
func (d *gLDriver) RenderWindowNow(win fyne.Window) bool {
	w, ok := win.(*window)
	if !ok || w.driver != d || w.closing || !w.visible || w.viewport == nil || nativeFrameDrawing {
		return false
	}
	if !w.canvas.CheckDirtyAndClear() {
		return false
	}
	nativeFrameDrawing = true
	defer func() { nativeFrameDrawing = false }()
	w.RunWithContext(func() { d.repaintWindow(w) })
	return true
}
