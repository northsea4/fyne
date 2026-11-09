//go:build !no_glfw && !mobile

package glfw

import (
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestWindow_CloseReleasesOnlyItsTextTexturesOnce(t *testing.T) {
	first, second := createWindow("first"), createWindow("second")
	firstKey := cache.FontCacheEntry{Canvas: first.window.canvas}
	secondKey := cache.FontCacheEntry{Canvas: second.window.canvas}
	firstFreed, secondFreed := 0, 0
	runOnMain(func() {
		cache.SetTextTexture(firstKey, cache.NoTexture, first.window.canvas, func() { firstFreed++ })
		cache.SetTextTexture(secondKey, cache.NoTexture, second.window.canvas, func() { secondFreed++ })
	})
	first.Close()
	runOnMain(func() {
		_, foundFirst := cache.GetTextTexture(firstKey)
		_, foundSecond := cache.GetTextTexture(secondKey)
		assert.False(t, foundFirst)
		assert.True(t, foundSecond)
		assert.Equal(t, 1, firstFreed)
		assert.Zero(t, secondFreed)
	})
	first.Close()
	second.Close()
	runOnMain(func() {
		assert.Equal(t, 1, firstFreed)
		assert.Equal(t, 1, secondFreed)
	})
}

func TestWindow_CloseReleasesCachedHiddenWidgets(t *testing.T) {
	w := createWindow("hidden cleanup")
	hidden := widget.NewButton("hidden action", nil)
	hidden.Hide()
	content := container.NewStack(widget.NewLabel("visible"), hidden)
	w.SetContent(content)
	unused := widget.NewButton("never rendered", nil)
	unused.Hide()
	runOnMain(func() {
		hidden.MinSize()
		assert.True(t, cache.IsRendered(hidden))
		content.Objects = append(content.Objects, unused)
		assert.False(t, cache.IsRendered(unused))
	})
	w.Close()
	runOnMain(func() {
		assert.False(t, cache.IsRendered(hidden))
		assert.False(t, cache.IsRendered(unused))
	})
}

func TestWindow_CloseDoesNotRebuildRenderersForPendingPointer(t *testing.T) {
	w := createWindow("closed pointer")
	label := widget.NewLabel("closed result")
	w.SetContent(label)
	w.Close()
	runOnMain(func() {
		assert.False(t, cache.IsRendered(label))
		w.window.mousePosUpdateProcessed = false
		w.window.newMousePosX, w.window.newMousePosY = 10, 10
		w.window.ensurePositionProcessed()
		assert.False(t, cache.IsRendered(label))
	})
}
