package gui

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
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

			editors := object.NewQuery().Where(func(c object.Component) bool {
				_, ok := c.(editor.T)
				return ok
			}).Collect(scene)

			fmt.Println("found", len(editors), "editors")

			for _, cmp := range editors {
				editor := cmp.(editor.T)
				editor.SelectColor(clr)
			}
		},
	})
}

/*
type Node interface {
	Render()
	Props() any
	Children() []Node
}

type node[T any] struct {
	element func(T)
	props T
	children []Node
}

func (n node[T]) Render() {
	n.element(n.props)
}

func (n node[T]) Props() any {
	return n.props
}

func (n node[T]) Children() []Node {
	return n.children
}

func CreateElement[T any](comp func(T), props T, children ...Node) Node {
	return node[T] {
		element: comp,
		props: props,
		children: children,
	}
}
*/
