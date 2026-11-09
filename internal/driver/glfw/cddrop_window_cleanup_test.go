//go:build !no_glfw && !mobile

package glfw

import (
	"testing"

	"fyne.io/fyne/v2/internal/cache"
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
