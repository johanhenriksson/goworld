package gui_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/gui/label"
	. "github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/palette"
	"github.com/johanhenriksson/goworld/gui/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

func TestReconcile1(t *testing.T) {
	a := palette.New("palette", &palette.Props{
		Palette: color.DefaultPalette,
	})
	b := palette.New("palette", &palette.Props{
		Palette: color.DefaultPalette,
	})
	Reconcile(a, b)
}
