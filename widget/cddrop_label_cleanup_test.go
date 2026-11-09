package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestLabelDestroyReleasesPrivateNonSelectableHelpers(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)
	label := NewLabel("result title")
	label.MinSize()
	label.Resize(fyne.NewSize(200, 30))
	selection, provider := label.selection, label.provider
	assert.True(t, cache.IsRendered(selection))
	assert.True(t, cache.IsRendered(provider))
	cache.DestroyRenderer(label)
	assert.False(t, cache.IsRendered(selection))
	assert.False(t, cache.IsRendered(provider))
	label.Refresh()
	label.Resize(fyne.NewSize(210, 30))
	assert.True(t, cache.IsRendered(label))
	assert.False(t, cache.IsRendered(selection))
	assert.False(t, cache.IsRendered(provider))
	cache.DestroyRenderer(label)
}
