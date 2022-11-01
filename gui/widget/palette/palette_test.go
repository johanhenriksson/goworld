package palette_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
)

func TestClickSwatch(t *testing.T) {
	app := node.NewRenderer(func() node.T {
		return palette.New("palette", palette.Props{
			Palette: color.DefaultPalette,
		})
	})
	view := vec2.New(1000, 1000)
	w := app.Render(view)

	swatch := widget.Find(w, "color1")
	if swatch == nil {
		t.Error("could not find swatch widget")
	}
	if swatch.Size() != vec2.New(195, 20) {
		t.Errorf("wrong swatch size: %s", swatch.Size())
	}

	// click color swatch
	widget.SimulateClick(swatch, mouse.Button1)

	// re-render
	w2 := app.Render(view)
	if w != w2 {
		t.Error("unexpected element recreation")
	}

	// lookup preview element
	preview := widget.Find(w2, "preview")
	if preview == nil {
		t.Error("could not find swatch widget")
	}

	// compare colors
	sp := swatch.Props().(rect.Props)
	pp := preview.Props().(rect.Props)
	if pp.Style.Color != sp.Style.Color {
		t.Error("expected preview color to be updated")
	}
}
