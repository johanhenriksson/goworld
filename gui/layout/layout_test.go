package layout_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/gui/dimension"
	. "github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"
)

func TestLayoutAbsolute(t *testing.T) {
	elsize := vec2.New(100, 100)
	rect := rect.Create("test", &rect.Props{
		Layout: Absolute{},
		Width:  dimension.Fixed(elsize.X),
		Height: dimension.Fixed(elsize.Y),
	})
	rect.Arrange(elsize)

	if rect.Arrange(elsize) != elsize {
		t.Errorf("wrong arrangement")
	}
}

func assertDimensions(t *testing.T, w widget.T, pos, size vec2.T) {
	t.Helper()
	if !w.Position().ApproxEqual(pos) {
		t.Errorf("expected widget position %s, got %s", pos, w.Position())
	}
	if !w.Size().ApproxEqual(size) {
		t.Errorf("expected widget size %s, got %s", size, w.Size())
	}
}
