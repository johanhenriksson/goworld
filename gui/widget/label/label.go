package label

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/kjk/flex"
)

type T interface {
	widget.T

	Text() string
}

type Props struct {
	Style   style.Sheet
	Text    string
	OnClick mouse.Callback
}

type label struct {
	widget.T
	Renderer

	props Props
	size  vec2.T
}

func New(key string, props Props) node.T {
	return node.Builtin(key, props, nil, new)
}

func new(key string, props Props) T {
	lbl := &label{
		T:        widget.New(key),
		Renderer: NewRenderer(),
	}
	lbl.Update(props)
	return lbl
}

func (l *label) Size() vec2.T { return l.T.Size() }

func (l *label) Props() any { return l.props }
func (l *label) Update(props any) {
	new := props.(Props)
	textChanged := new.Text != l.props.Text
	styleChanged := new.Style != l.props.Style
	l.props = new

	if styleChanged {
		new.Style.Apply(l)
		l.Flex().MarkDirty()
	}

	if textChanged {
		l.Renderer.SetText(new.Text)
		l.Flex().MarkDirty()
	}
}

// prop accessors

func (l *label) Text() string { return l.props.Text }

func (l *label) Draw(args render.Args) {
	l.T.Draw(args)
	l.Renderer.Draw(args)
}

func (l *label) Flex() *flex.Node {
	node := l.T.Flex()
	node.SetMeasureFunc(l.Renderer.Measure)
	return node
}

func (l *label) Destroy() {
	l.Renderer.Destroy()
}

//
// Events
//

func (l *label) MouseEvent(e mouse.Event) {
	if e.Action() == mouse.Press && l.props.OnClick != nil {
		l.props.OnClick(e)
		e.Consume()
	}
}
