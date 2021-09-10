package widget

import (
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

type Props interface{}

type T interface {
	Key() string
	Props() Props
	Update(Props)
	Destroy()

	Size() vec2.T
	Resize(vec2.T)
	Position() vec2.T
	Move(vec2.T)
	Width() dimension.T
	Height() dimension.T

	Reflow()
	Draw(render.Args)
}

type widget struct {
	key      string
	size     vec2.T
	position vec2.T
	children []T
}

func New(key string, children ...T) T {
	return &widget{
		key:      key,
		children: children,
	}
}

func (w *widget) Key() string      { return w.key }
func (w *widget) Props() Props     { return nil }
func (w *widget) Position() vec2.T { return w.position }
func (w *widget) Size() vec2.T     { return w.size }
func (w *widget) Move(p vec2.T)    { w.position = p }
func (w *widget) Resize(s vec2.T)  { w.size = s }
func (w *widget) Draw(render.Args) {}
func (w *widget) Reflow()          {}

func (w *widget) Width() dimension.T {
	return dimension.Auto()
}

func (w *widget) Height() dimension.T {
	return dimension.Auto()
}

func (w *widget) Update(Props) {}
func (w *widget) Destroy()     {}
