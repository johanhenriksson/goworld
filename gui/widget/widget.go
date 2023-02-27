package widget

import (
	"fmt"

	"github.com/johanhenriksson/goworld/math/vec2"

	"github.com/kjk/flex"
)

var FlexConfig = flex.NewConfig()

type T interface {
	Key() string

	// Properties returns a pointer to the components property struct.
	// The pointer is used to compare the states when deciding if the component needs to be updated.
	Props() any

	// Update replaces the components property struct.
	Update(any)

	// Destroy releases any resources associated with the component.
	// Attempting to draw a destroyed component will cause a panic.
	Destroy()

	// Destroyed indicates whether the component has been destroyed or not.
	Destroyed() bool

	// Size returns the actual size of the element in pixels
	Size() vec2.T

	// Position returns the current position of the element relative to its parent
	Position() vec2.T

	ZOffset() int

	Children() []T
	SetChildren([]T)

	Flex() *flex.Node

	// Draw the widget. This should only be called by the GUI Draw Pass
	// Calling Draw() will instantiate any required GPU resources prior to drawing.
	// Attempting to draw a destroyed component will cause a panic.
	Draw(DrawArgs)

	Draw2(DrawArgs, *QuadBuffer)
}

type widget struct {
	key       string
	destroyed bool
	flex      *flex.Node
}

func New(key string) T {
	node := flex.NewNodeWithConfig(FlexConfig)
	node.Context = key
	return &widget{
		key:  key,
		flex: node,
	}
}

func (w *widget) Key() string     { return w.key }
func (w *widget) Destroyed() bool { return w.destroyed }
func (w *widget) Destroy()        { w.destroyed = true }

func (w *widget) Flex() *flex.Node {
	return w.flex
}

func (w *widget) Position() vec2.T {
	//return vec2.New(w.flex.LayoutGetLeft(), w.flex.LayoutGetTop())
	return calculatePosition(w.flex)
}

func (w *widget) ZOffset() int {
	return 0
}

func (w *widget) Size() vec2.T {
	return vec2.New(w.flex.LayoutGetWidth(), w.flex.LayoutGetHeight())
}

func (w *widget) Props() any {
	panic("widget.Props() must be implemented")
}

func (w *widget) Update(any) {
	panic("widget.Update() must be implemented")
}

func (w *widget) Children() []T { return nil }
func (w *widget) SetChildren(c []T) {
	if len(c) > 0 {
		panic("widget cant have children")
	}
}

func (w *widget) Draw(DrawArgs) {
	// base widget Draw() should be called ahead of overridden draws

	if w.Destroyed() {
		panic(fmt.Sprintf("attempt to draw destroyed widget %s", w.key))
	}
}

func (w *widget) Draw2(DrawArgs, *QuadBuffer) {
}

func calculatePosition(node *flex.Node) vec2.T {
	pos := vec2.New(node.LayoutGetLeft(), node.LayoutGetTop())
	if node.Parent != nil {
		pos = pos.Add(calculatePosition(node.Parent))
	}
	return pos
}
