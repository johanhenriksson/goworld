package component

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Props interface{}

type T interface {
	widget.T
}

type component[P any] struct {
	widget.T
	props    P
	children []widget.T
}

func Create[P Props](w widget.T, props P) T {
	return &component[P]{
		T:     w,
		props: props,
	}
}

func (f *component[P]) Draw(args widget.DrawArgs) {
	for _, child := range f.Children() {
		// calculate child transform
		// try to fix the position to an actual pixel
		// pos := vec3.Extend(child.Position().Scaled(args.Viewport.Scale).Floor().Scaled(1/args.Viewport.Scale), -1)
		z := child.ZOffset()
		pos := vec3.Extend(child.Position(), args.Position.Z-float32(1+z))
		transform := mat4.Translate(pos)
		childArgs := args
		childArgs.Transform = transform // .Mul(&args.Transform)
		childArgs.Position = pos

		// draw child
		child.Draw(childArgs)
	}
}

func (f *component[P]) ZOffset() int {
	return 0
}

func (f *component[P]) Children() []widget.T { return f.children }
func (f *component[P]) SetChildren(c []widget.T) {
	// reconcile flex children
	if len(c) > 1 {
		panic("expected component to have a single child")
	}
	f.children = c
	if len(c) == 0 {
		return
	}
	child := c[0]

	if len(f.Flex().Children) == 1 {
		if f.Flex().Children[0] == child.Flex() {
			return
		}
		// replace
		f.Flex().RemoveChild(f.Flex().Children[0])
	}
	f.Flex().InsertChild(child.Flex(), 0)
}

//
// Lifecycle
//

func (f *component[P]) Props() any { return f.props }

func (f *component[P]) Update(p any) {
	var ok bool
	f.props, ok = p.(P)
	if !ok {
		panic("illegal props")
	}
}

func (f *component[P]) Destroy() {
	f.T.Destroy()

	for _, child := range f.children {
		child.Destroy()
	}
}

//
// Events
//

func (f *component[P]) MouseEvent(e mouse.Event) {
	// because children may have absolute positioning, we must pass the event to all of them.
	// children always have higher z index, so they have priority
	for _, frame := range f.children {
		if handler, ok := frame.(mouse.Handler); ok {
			handler.MouseEvent(e)
			if e.Handled() {
				e.Consume()
				return
			}
		}
	}
}
