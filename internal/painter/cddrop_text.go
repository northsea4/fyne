package painter

import (
	"math"

	"fyne.io/fyne/v2"
	"github.com/go-text/typesetting/shaping"
)

// CDDropTextInkBounds measures visible glyphs using the very same shaped runs
// and font cache as DrawString. It does not parse or retain another CJK font.
// Coordinates are relative to the normal Fyne text origin, before line centering.
func CDDropTextInkBounds(text string, size float32, style fyne.TextStyle, source fyne.Resource) (top, bottom float32) {
	face := CachedFontFace(style, source, nil)
	top, bottom = float32(math.Inf(1)), float32(math.Inf(-1))
	advance := float32(0)
	walkString(face.Fonts, text, float32ToFixed266(size), style, &advance, 1, func(run shaping.Output, _, baseline float32) {
		for _, glyph := range run.Glyphs {
			if glyph.Width == 0 || glyph.Height == 0 {
				continue
			}
			y := baseline - fixed266ToFloat32(glyph.YBearing+glyph.YOffset)
			end := y - fixed266ToFloat32(glyph.Height)
			top = min(top, y)
			bottom = max(bottom, end)
		}
	})
	if math.IsInf(float64(top), 0) {
		top, bottom = 0, size
	}
	return top, bottom
}
