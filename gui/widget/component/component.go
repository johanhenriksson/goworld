package component

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"

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
	flex     *flex.Node
}

func New(key string, props any) T {
	return &component{
		key:   key,
		props: props,
		flex:  flex.NewNodeWithConfig(widget.FlexConfig),
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
	// whats the reasoning here?
	// why do we pick the first element
	if len(children) > 0 {
		c.wrap = children[0]
	} else {
		c.wrap = nil
	}
}

func (c *component) Draw(args widget.DrawArgs) {
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

func (c *component) ZOffset() int {
	if c.wrap != nil {
		return c.wrap.ZOffset()
	}
	return 0
}

func (c *component) Flex() *flex.Node {
	if c.wrap != nil {
		return c.wrap.Flex()
	}
	return c.flex
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
			handler.MouseEvent(e)
			if e.Handled() {
				return
			}
		}
	}
}
