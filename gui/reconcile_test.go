package gui_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	. "github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
)

func findWidget(key string, w widget.T) widget.T {
	if w.Key() == key {
		return w
	}
	for _, child := range w.Children() {
		if hit := findWidget(key, child); hit != nil {
			return hit
		}
	}
	return nil
}

func TestReconcile1(t *testing.T) {
	a := palette.New("palette", &palette.Props{
		Palette: color.DefaultPalette,
	})
	a.Render(nil)
	w := a.Hydrate()

	swatch := findWidget("color1", w).(rect.T)
	if swatch == nil {
		t.Error("could not find swatch widget")
	}

	sp := swatch.Props().(*rect.Props)

	// simulate mouse click
	sp.OnClick(mouse.NewButtonEvent(mouse.Button1, mouse.Press, vec2.New(0, 0), 0, false))

	b := palette.New("palette", &palette.Props{
		Palette: color.DefaultPalette,
	})

	r := Reconcile(a, b)
	w2 := r.Hydrate()
	preview := findWidget("preview", w2)
	if preview == nil {
		t.Error("could not find swatch widget")
	}

	pp := preview.Props().(*rect.Props)
	if pp.Color != sp.Color {
		t.Error("expected preview color to be updated")
	}
}
