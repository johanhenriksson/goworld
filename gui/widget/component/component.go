package component

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/kjk/flex"
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

func (c *component) Update(props any) {
	c.props = props
}

func (c *component) Props() any {
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

func (c *component) Size() vec2.T {
	if c.wrap != nil {
		return c.wrap.Size()
	}
	return vec2.Zero
}

func (c *component) Position() vec2.T {
	if c.wrap != nil {
		return c.wrap.Position()
	}
	return vec2.Zero
}

func (c *component) Style() style.Sheet {
	if c.wrap != nil {
		return c.wrap.Style()
	}
	return style.Sheet{}
}

func (c *component) SetStyle(style style.Sheet) {
	if c.wrap != nil {
		c.wrap.SetStyle(style)
	}
}

func (c *component) Flex() *flex.Node {
	if c.wrap != nil {
		return c.wrap.Flex()
	}
	return nil
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
