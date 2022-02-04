package gui

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/label"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/palette"
	"github.com/johanhenriksson/goworld/render/color"
)

func CounterLabel(key, format string) node.T {
	count, setCount := hooks.UseState[int](0)

	return label.New(key, &label.Props{
		Text:  fmt.Sprintf(format, count),
		Size:  16.0,
		Color: color.White,
		OnClick: func(e mouse.Event) {
			setCount(count + 1)
		},
	})
}

func pickColor(clr color.T) {
	scene := hooks.UseScene()
	fmt.Println("pick callback:", clr)

	editor := query.New[editor.T]().First(scene)
	if editor == nil {
		fmt.Println("could not find editor")
		return
	}

	editor.SelectColor(clr)
}

func TestUI() node.T {
	return palette.New("palette", &palette.Props{
		Palette: color.DefaultPalette,
		OnPick:  pickColor,
	})
}
