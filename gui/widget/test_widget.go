package widget

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/kjk/flex"
)

func Dummy(key string) T {
	node := flex.NewNodeWithConfig(flex.NewConfig())
	node.Context = key
	return &dummy{
		key:  key,
		flex: node,
	}
}

type dummy struct {
	key      string
	flex     *flex.Node
	children []T
	props    any
}

func (w *dummy) Key() string { return w.key }
func (w *dummy) Destroy()    {}

func (w *dummy) Update(p any) { w.props = p }
func (w *dummy) Props() any   { return w.props }

func (w *dummy) Flex() *flex.Node { return w.flex }

func (w *dummy) Position() vec2.T { return vec2.Zero }
func (w *dummy) Size() vec2.T     { return vec2.Zero }

func (w *dummy) Draw(args DrawArgs, quads *QuadBuffer) {}

func (w *dummy) Children() []T     { return w.children }
func (w *dummy) SetChildren(c []T) { w.children = c }
