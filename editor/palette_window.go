package editor

import (
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/ui"
)

type PaletteWindow struct {
	*ui.Rect
	Palette  render.Palette
	Selected render.Color
}

func NewPaletteWindow(palette render.Palette) *PaletteWindow {
	cols := 5
	gridStyle := ui.Style{"layout": ui.String("column"), "spacing": ui.Float(2)}
	rowStyle := ui.Style{"layout": ui.String("row"), "spacing": ui.Float(2)}
	rows := make([]ui.Component, 0, len(palette)/cols+1)
	row := make([]ui.Component, 0, cols)

	wnd := &PaletteWindow{
		Palette:  palette,
		Selected: palette[0],
	}

	for i := 1; i <= len(palette); i++ {
		itemIdx := i - 1
		color := palette[itemIdx]

		swatch := ui.NewRect(ui.Style{"color": ui.Color(color), "layout": ui.String("fixed")})
		swatch.Resize(vec2.New(20, 20))
		swatch.OnClick(func(ev ui.MouseEvent) {
			if ev.Button == mouse.Button1 {
				wnd.Selected = color
			}
		})

		row = append(row, swatch)

		if i%cols == 0 {
			rows = append(rows, ui.NewRect(rowStyle, row...))
			row = make([]ui.Component, 0, cols)
		}
	}

	wnd.Rect = ui.NewRect(WindowStyle,
		ui.NewText("Palette", ui.NoStyle),
		ui.NewRect(gridStyle, rows...))

	wnd.Flow(vec2.New(200, 400))

	return wnd
}
