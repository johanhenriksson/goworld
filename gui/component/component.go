package component

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

type T interface {
	widget.T
}

type component struct {
	widget.T
	props    any
	children []widget.T
}

func New(key string, props any) widget.T {
	return &component{
		T:     widget.New(key),
		props: props,
	}
}

func (c *component) Update(props widget.Props) {
	c.props = props
}

func (c *component) Children() []widget.T {
	return c.children
}

func (c *component) SetChildren(children []widget.T) {
	c.children = children
}

func (c *component) Draw(args render.Args) {
	for _, child := range c.children {
		child.Draw(args)
	}
}
func (c *component) Reflow() {
	for _, child := range c.children {
		child.Reflow()
	}
}
func (c *component) Resize(s vec2.T) {
	for _, child := range c.children {
		child.Resize(s)
	}
}

func (c *component) MouseEvent(e mouse.Event) {
	fmt.Println("component mouse event")
	for _, frame := range c.children {
		if handler, ok := frame.(mouse.Handler); ok {
			ev := e.Project(frame.Position())
			target := ev.Position()
			size := frame.Size()
			if target.X < 0 || target.X > size.X || target.Y < 0 || target.Y > size.Y {
				// outside
				continue
			}

			handler.MouseEvent(ev)
			if ev.Handled() {
				e.Consume()
				return
			}
		}
	}
}
