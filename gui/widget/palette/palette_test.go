package palette_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/palette"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

func TestClickSwatch(t *testing.T) {
	app := node.NewRenderer(func() node.T {
		return palette.New("palette", &palette.Props{
			Palette: color.DefaultPalette,
		})
	})
	w := app.Render()

	swatch := widget.Find(w, "color1")
	if swatch == nil {
		t.Error("could not find swatch widget")
	}

	// click color swatch
	widget.SimulateClick(swatch, mouse.Button1)

	// re-render
	w2 := app.Render()
	if w != w2 {
		t.Error("unexpected element recreation")
	}

	// lookup preview element
	preview := widget.Find(w2, "preview")
	if preview == nil {
		t.Error("could not find swatch widget")
	}

	// compare colors
	sp := swatch.Props().(*rect.Props)
	pp := preview.Props().(*rect.Props)
	if pp.Color != sp.Color {
		t.Error("expected preview color to be updated")
	}
}
