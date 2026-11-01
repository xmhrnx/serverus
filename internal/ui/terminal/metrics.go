package terminal

import (
	"math"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"golang.org/x/image/math/fixed"

	"github.com/mehran9675/giode/fonts"
)

// cellMetrics is the size of one character cell in pixels.
type cellMetrics struct {
	width    float32 // advance of one monospaced character
	height   int     // line height
	baseline int     // offset of the text's top inside the line
	key      [2]float32
}

// lineSpacing adds air between lines, as terminals usually do.
const lineSpacing = 1.22

// measure returns the cell size for the mono face at size, reusing the
// last result while the size and display scale are unchanged.
func (m *cellMetrics) measure(gtx layout.Context, face font.Font, size unit.Sp) {
	key := [2]float32{float32(size), gtx.Metric.PxPerSp}
	if m.key == key && m.width > 0 {
		return
	}
	sh := fonts.Shaper()
	px := gtx.Sp(size)
	m.key = key
	m.width = float32(px) * 0.6
	m.height = int(math.Ceil(float64(px) * lineSpacing))
	if sh == nil {
		return
	}
	sh.LayoutString(text.Parameters{Font: face, PxPerEm: fixed.I(px), MaxLines: 1}, "M")
	if g, ok := sh.NextGlyph(); ok {
		m.width = float32(g.Advance) / 64
		glyphHeight := float64(g.Ascent+g.Descent) / 64
		m.height = int(math.Ceil(glyphHeight * lineSpacing))
		m.baseline = int((float64(m.height) - glyphHeight) / 2)
	}
	for {
		if _, ok := sh.NextGlyph(); !ok {
			break
		}
	}
}
