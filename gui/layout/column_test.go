package layout_test

import (
	. "github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/rect"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"

	"testing"
)

func TestColumnLayout(t *testing.T) {
	a := widget.New("a")
	b := widget.New("b")
	parent := rect.New("parent", &rect.Props{}, a, b)
	parent.Resize(vec2.New(120, 70))

	c := Column{
		Padding: 10,
		Gutter:  10,
	}
	c.Flow(parent)

	// inner bounds should be 100x50
	// gutter is 10px -> usable space 100x40

	assertDimensions(t, a, vec2.New(10, 10), vec2.New(100, 20))
	assertDimensions(t, b, vec2.New(10, 40), vec2.New(100, 20))
}

func assertDimensions(t *testing.T, w widget.T, pos, size vec2.T) {
	t.Helper()
	if !w.Position().ApproxEqual(pos) {
		t.Errorf("expected widget position %s, got %s", pos, w.Position())
	}
	if !w.Size().ApproxEqual(size) {
		t.Errorf("expected widget size %s, got %s", pos, w.Size())
	}
}
