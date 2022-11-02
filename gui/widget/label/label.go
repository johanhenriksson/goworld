package label

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/kjk/flex"
)

type ChangeCallback func(string)

type T interface {
	widget.T
	style.FontWidget

	Text() string
}

type Props struct {
	Style    Style
	Text     string
	OnClick  mouse.Callback
	OnChange ChangeCallback
}

type label struct {
	widget.T
	Renderer

	props Props
	scale float32
	state style.State
}

func New(key string, props Props) node.T {
	return node.Builtin(key, props, nil, new)
}

func new(key string, props Props) T {
	lbl := &label{
		T:        widget.New(key),
		Renderer: NewRenderer(),
		scale:    1,
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
		new.Style.Apply(l, l.state)
		l.Flex().MarkDirty()
	}

	if textChanged {
		l.Renderer.SetText(new.Text)
		l.Flex().MarkDirty()
	}
}

// prop accessors

func (l *label) Text() string { return l.props.Text }

func (l *label) Draw(args widget.DrawArgs) {
	if l.props.Style.Hidden {
		return
	}

	if window.Scale != l.scale {
		// ui scale has changed
		l.Flex().MarkDirty()
		l.scale = window.Scale
	}

	l.T.Draw(args)
	l.Renderer.Draw(args, l)
}

func (l *label) Flex() *flex.Node {
	node := l.T.Flex()
	node.SetMeasureFunc(l.Renderer.Measure)
	return node
}

//
// Events
//

func (l *label) MouseEvent(e mouse.Event) {
	target := e.Position().Sub(l.Position())
	size := l.Size()
	mouseover := target.X >= 0 && target.X < size.X && target.Y >= 0 && target.Y < size.Y

	if mouseover {
		// hover start
		if !l.state.Hovered {
			l.state.Hovered = true
			l.props.Style.Apply(l, l.state)
		}

		// click event
		// todo: separate into mouse down/up?
		if e.Action() == mouse.Press && l.props.OnClick != nil {
			l.props.OnClick(e)
			e.Consume()
		}

		// take input keyboard focus
		if e.Action() == mouse.Press || e.Action() == mouse.Release {
			if l.props.OnChange != nil {
				keys.Focus(l)
				e.Consume()
			}
		}
	} else {
		// hover end
		if l.state.Hovered {
			l.state.Hovered = false
			l.props.Style.Apply(l, l.state)
		}
	}
}

func (l *label) KeyEvent(e keys.Event) {
	if l.props.OnChange == nil {
		return
	}
	if e.Action() == keys.Char {
		l.props.OnChange(l.props.Text + string(e.Character()))
	}
	if e.Action() == keys.Press || e.Action() == keys.Repeat {
		switch e.Code() {
		case keys.Backspace:
			l.props.OnChange(l.props.Text[:len(l.props.Text)-1])
		}
	}
}
