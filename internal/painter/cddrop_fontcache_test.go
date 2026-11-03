//go:build go1.24

package painter

import (
	"runtime"
	"testing"
	"weak"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/go-text/typesetting/font"
)

func shapedFontReference() weak.Pointer[font.Font] {
	resource := fyne.NewStaticResource("cache-lifecycle.ttf", theme.DefaultTextFont().Content())
	face := CachedFontFace(fyne.TextStyle{}, resource, nil)
	MeasureString(face.Fonts, "Theme changes", 14, fyne.TextStyle{})
	return weak.Make(face.Fonts.ResolveFace('T').Font)
}

// Clearing the theme font caches must also drop fonts held by the shaper's LRU.
func TestCDDropThemeClearReleasesObsoleteFont(t *testing.T) {
	ClearFontCache()
	reference := shapedFontReference()
	runtime.GC()
	if reference.Value() == nil {
		t.Fatal("font was not cached before clear")
	}
	ClearFontCache()
	runtime.GC()
	if reference.Value() != nil {
		t.Fatal("obsolete parsed font retained after theme cache clear")
	}
}
