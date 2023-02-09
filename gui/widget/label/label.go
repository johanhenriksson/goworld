package label

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"

	"github.com/kjk/flex"
	"golang.org/x/exp/utf8string"
)

type ChangeCallback func(string)

type T interface {
	widget.T
	style.FontWidget
	keys.Handler

	Text() string
	Cursor() int
}

type Props struct {
	Style    Style
	Text     string
	OnClick  mouse.Callback
	OnChange ChangeCallback

	OnKeyUp   keys.Callback
	OnKeyDown keys.Callback
	OnKeyChar keys.Callback
}

type label struct {
	widget.T
	Renderer

	props  Props
	scale  float32
	state  style.State
	cursor int
	text   string
}

func New(key string, props Props) node.T {
	return node.Builtin(key, props, nil, new)
}

func new(w widget.T, props Props) T {
	lbl := &label{
		T:        w,
		Renderer: NewRenderer(w.Key()),
		scale:    1,
		cursor:   utf8string.NewString(props.Text).RuneCount(),
		text:     props.Text,
		props: Props{
			Text: "\x00",
		},
	}
	lbl.Update(props)
	return lbl
}

func (l *label) Size() vec2.T   { return l.T.Size() }
func (l *label) editable() bool { return l.props.OnChange != nil }

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
		l.setText(new.Text)
		l.Flex().MarkDirty()
	}
}

func (l *label) setText(text string) {
	str := utf8string.NewString(text)

	if text != l.text {
		// new text is different from what we had before,
		// move cursor to end of line
		l.cursor = str.RuneCount()
	}

	l.cursor = math.Min(l.cursor, str.RuneCount())
	l.text = text

	if l.editable() {
		// we also need to know if it has focus
		text = str.Slice(0, l.cursor) + "_" + str.Slice(l.cursor, str.RuneCount())
		// text = l.props.Text[:l.cursor] + "_" + l.props.Text[l.cursor:]
	}

	// update renderer
	l.Renderer.SetText(text)
}

// prop accessors

func (l *label) Text() string { return l.props.Text }
func (l *label) Cursor() int  { return l.cursor }

func (l *label) Draw(args widget.DrawArgs) {
	if l.props.Style.Hidden {
		return
	}

	if args.Viewport.Scale != l.scale {
		// ui scale has changed
		l.Flex().MarkDirty()
		l.scale = args.Viewport.Scale
	}

	l.T.Draw(args)
	l.Renderer.Draw(args, l)
}

func (l *label) Flex() *flex.Node {
	node := l.T.Flex()
	node.SetMeasureFunc(l.measure)
	return node
}

func (l *label) measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size {
	r := l.Renderer.Measure(node, width, widthMode, height, heightMode)

	if widthMode == flex.MeasureModeExactly {
		r.Width = width
	}
	if widthMode == flex.MeasureModeAtMost {
		r.Width = math.Min(r.Width, width)
	}

	return r
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
			if l.editable() {
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
	//
	// key events
	//

	if l.props.OnKeyUp != nil && e.Action() == keys.Release {
		l.props.OnKeyUp(e)
		if e.Handled() {
			return
		}
	}
	if l.props.OnKeyDown != nil && e.Action() == keys.Press {
		l.props.OnKeyDown(e)
		if e.Handled() {
			return
		}
	}
	if l.props.OnKeyChar != nil && e.Action() == keys.Char {
		l.props.OnKeyChar(e)
		if e.Handled() {
			return
		}
	}

	//
	// text state handling
	//

	if l.props.OnChange == nil {
		return
	}
	if e.Action() == keys.Char {
		str := utf8string.NewString(l.text)
		l.text = str.Slice(0, l.cursor) + string(e.Character()) + str.Slice(l.cursor, str.RuneCount())
		l.cursor = math.Min(l.cursor+1, utf8string.NewString(l.text).RuneCount())
		l.props.OnChange(l.text)
	}
	if e.Action() == keys.Press || e.Action() == keys.Repeat {
		switch e.Code() {
		case keys.Backspace:
			str := utf8string.NewString(l.text)
			if l.cursor > 0 {
				l.cursor--
				l.text = str.Slice(0, l.cursor) + str.Slice(l.cursor+1, str.RuneCount())
				l.props.OnChange(l.text)
			}

		case keys.Delete:
			str := utf8string.NewString(l.text)
			if l.cursor < str.RuneCount() {
				l.text = str.Slice(0, l.cursor) + str.Slice(l.cursor+1, str.RuneCount())
				l.props.OnChange(l.text)
			}

		case keys.LeftArrow:
			l.cursor = math.Clamp(l.cursor-1, 0, len(l.text))
			l.setText(l.text)

		case keys.RightArrow:
			l.cursor = math.Clamp(l.cursor+1, 0, len(l.text))
			l.setText(l.text)

		case keys.U:
			// ctrl+u clears text
			if e.Modifier(keys.Ctrl) {
				l.setText("")
				l.props.OnChange("")
			}
		}
	}
}
