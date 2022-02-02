package gui

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/label"
	"github.com/johanhenriksson/goworld/gui/palette"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/render/color"
)

func CounterLabel(key, format string) widget.T {
	count, setCount := hooks.UseInt(0)

	return label.New(key, &label.Props{
		Text:  fmt.Sprintf(format, count),
		Size:  16.0,
		Color: color.White,
		OnClick: func(e mouse.Event) {
			setCount(count + 1)
		},
	})
}

func TestUI() widget.T {
	scene := hooks.UseScene()
	return palette.New("palette", &palette.Props{
		Palette: color.DefaultPalette,
		OnPick: func(clr color.T) {
			fmt.Println("pick callback:", clr)

			editor := query.New[editor.T]().First(scene)
			if editor == nil {
				fmt.Println("could not find editor")
				return
			}

			editor.SelectColor(clr)
		},
	})
}
