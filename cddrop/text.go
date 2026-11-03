// Package cddrop exposes the app's small text-metric extension to Fyne.
package cddrop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/painter"
)

// TextInkBounds reuses Fyne's shaping cache; it measures visible ink, not the
// font's larger ascender/descender line box. See ../CDDROP_PATCH.md.
func TextInkBounds(text string, size float32, style fyne.TextStyle, source fyne.Resource) (top, bottom float32) {
	return painter.CDDropTextInkBounds(text, size, style, source)
}
