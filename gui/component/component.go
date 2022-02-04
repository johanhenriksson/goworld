package component

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

type T interface {
	widget.T
}

type component struct {
	key      string
	wrap     widget.T
	children []widget.T
	props    any
}

func New(key string, props any) T {
	return &component{
		key:   key,
		props: props,
	}
}

func (c *component) Key() string {
	return c.key
}

func (c *component) Update(props widget.Props) {
	c.props = props
}

func (c *component) Props() widget.Props {
	return c.props
}

func (c *component) Children() []widget.T {
	return c.children
}

func (c *component) SetChildren(children []widget.T) {
	c.children = children
	if len(children) > 0 {
		c.wrap = children[0]
	} else {
		c.wrap = nil
	}
}

func (c *component) Draw(args render.Args) {
	if c.wrap != nil {
		c.wrap.Draw(args)
	}
}

func (c *component) Reflow() {
	if c.wrap != nil {
		c.wrap.Reflow()
	}
}
func (c *component) Resize(s vec2.T) {
	if c.wrap != nil {
		c.wrap.Resize(s)
	}
}

func (c *component) Size() vec2.T {
	if c.wrap != nil {
		return c.wrap.Size()
	}
	return vec2.Zero
}

func (c *component) Move(t vec2.T) {
	if c.wrap != nil {
		c.wrap.Move(t)
	}
}
func (c *component) Position() vec2.T {
	if c.wrap != nil {
		return c.wrap.Size()
	}
	return vec2.Zero
}

func (c *component) Width() dimension.T {
	if c.wrap != nil {
		return c.wrap.Width()
	}
	return dimension.Auto()
}

func (c *component) Height() dimension.T {
	if c.wrap != nil {
		return c.wrap.Height()
	}
	return dimension.Auto()
}

func (c *component) Destroy() {
	if c.wrap != nil {
		c.wrap.Destroy()
	}
}

func (c *component) Destroyed() bool {
	if c.wrap != nil {
		return c.wrap.Destroyed()
	}
	return false
}

func (c *component) MouseEvent(e mouse.Event) {
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