package rect_test

import (
	. "github.com/johanhenriksson/goworld/gui/rect"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"

	"testing"
)

func TestColumnLayout(t *testing.T) {
	props := &Props{
		Padding: 10,
		Gutter:  10,
	}
	a := New("a", &Props{})
	b := New("b", &Props{})
	parent := New("parent", props, a, b)
	parent.Resize(vec2.New(120, 70))
	Column(parent, props)

	// inner bounds should be 100x50
	// gutter is 10px -> usable space 100x40

	assertDimensions(t, a, vec2.New(10, 10), vec2.New(100, 20))
	assertDimensions(t, b, vec2.New(10, 40), vec2.New(100, 20))
}

func TestRowLayout(t *testing.T) {
	props := &Props{
		Padding: 10,
		Gutter:  10,
	}
	a := New("a", &Props{})
	b := New("b", &Props{})
	parent := New("parent", props, a, b)
	parent.Resize(vec2.New(130, 60))
	Row(parent, props)

	// inner bounds should be 110x40
	// gutter is 10px -> usable space 100x40

	assertDimensions(t, a, vec2.New(10, 10), vec2.New(50, 40))
	assertDimensions(t, b, vec2.New(70, 10), vec2.New(50, 40))
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
