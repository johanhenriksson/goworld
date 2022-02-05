package layout_test

import (
	. "github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"

	"testing"
)

func TestRowLayout(t *testing.T) {
	a := widget.New("a")
	b := widget.New("b")
	parent := rect.New("parent", &rect.Props{}, a, b)
	parent.Resize(vec2.New(130, 60))

	r := Row{
		Padding: 10,
		Gutter:  10,
	}
	r.Flow(parent)

	// inner bounds should be 110x40
	// gutter is 10px -> usable space 100x40

	assertDimensions(t, a, vec2.New(10, 10), vec2.New(50, 40))
	assertDimensions(t, b, vec2.New(70, 10), vec2.New(50, 40))
}
