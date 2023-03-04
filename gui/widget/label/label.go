package label

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/font"
	"github.com/johanhenriksson/goworld/render/texture"

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
	OnBlur   func()

	OnKeyUp   keys.Callback
	OnKeyDown keys.Callback
	OnKeyChar keys.Callback
}

type label struct {
	key    string
	flex   *flex.Node
	props  Props
	scale  float32
	state  style.State
	cursor int
	tex    texture.Ref

	version    int
	text       string
	size       int
	fontName   string
	font       font.T
	color      color.T
	lineHeight float32
	texSize    vec2.T
}

func New(key string, props Props) node.T {
	return node.Builtin(key, props, nil, new)
}

func new(key string, props Props) T {
	node := flex.NewNodeWithConfig(flex.NewConfig())
	node.Context = key
	lbl := &label{
		key:    key,
		flex:   node,
		scale:  1,
		cursor: utf8string.NewString(props.Text).RuneCount(),
		text:   props.Text,

		lineHeight: 1,
		props: Props{
			Text: "\x00",
		},
	}
	lbl.flex.SetMeasureFunc(lbl.measure)
	lbl.Update(props)
	return lbl
}

//
// Widget implementation
//

func (l *label) Key() string          { return l.key }
func (l *label) Flex() *flex.Node     { return l.flex }
func (l *label) Position() vec2.T     { return vec2.New(l.flex.LayoutGetLeft(), l.flex.LayoutGetTop()) }
func (l *label) Size() vec2.T         { return vec2.New(l.flex.LayoutGetWidth(), l.flex.LayoutGetHeight()) }
func (l *label) Children() []widget.T { return nil }

func (l *label) SetChildren(children []widget.T) {
	if len(children) > 0 {
		panic("labels may not have children")
	}
}

func (l *label) Props() any { return l.props }
func (l *label) Update(props any) {
	new := props.(Props)
	textChanged := new.Text != l.props.Text
	styleChanged := new.Style != l.props.Style
	l.props = new

	if styleChanged {
		new.Style.Apply(l, l.state)
	}

	if textChanged {
		l.setText(new.Text)
	}
}

func (l *label) invalidate() {
	l.Flex().MarkDirty()
	if l.font == nil {
		fontName := l.fontName
		if fontName == "" {
			fontName = DefaultFont.Name
		}
		size := l.size
		if size == 0 {
			size = DefaultFont.Size
		}
		l.font = assets.GetFont(fontName, size, l.scale)
	}

	fargs := font.Args{
		LineHeight: l.lineHeight,
		Color:      color.White,
	}
	l.version++

	// todo: immediate updates causes a noticable flash while the text is rendered
	//       keep a reference until the new texture is ready?
	l.tex = font.Ref(l.Key(), l.version, l.font, l.text, fargs)
	l.texSize = l.font.Measure(l.text, fargs)
}

func (l *label) editable() bool { return l.props.OnChange != nil }

func (l *label) setText(text string) {
	if text == l.text {
		return
	}

	str := utf8string.NewString(text)

	l.cursor = math.Min(l.cursor, str.RuneCount())
	l.text = text

	// shit attempt at a cursor.
	// todo: make a real cursor. a blinking one
	// if l.editable() {
	// we also need to know if it has focus
	// text = str.Slice(0, l.cursor) + "_" + str.Slice(l.cursor, str.RuneCount())
	// text = l.props.Text[:l.cursor] + "_" + l.props.Text[l.cursor:]
	// }

	l.invalidate()
}

//
// Styles
//

func (l *label) SetFont(font style.Font) {
	if font.Name == l.fontName && font.Size == l.size {
		return
	}
	l.fontName = font.Name
	l.size = font.Size
	l.font = nil
	l.invalidate()
}

func (l *label) SetFontColor(color color.T) {
	l.color = color
}

func (l *label) SetLineHeight(lineHeight float32) {
	if lineHeight == l.lineHeight {
		return
	}
	l.lineHeight = lineHeight
	l.invalidate()
}

func (l *label) Text() string { return l.props.Text }
func (l *label) Cursor() int  { return l.cursor }

//
// Draw
//

func (l *label) Draw(args widget.DrawArgs, quads *widget.QuadBuffer) {
	if l.props.Style.Hidden {
		return
	}

	// skip empty strings
	if l.text == "" {
		return
	}

	if args.Viewport.Scale != l.scale {
		// ui scale has changed
		l.scale = args.Viewport.Scale
		l.font = nil
		l.invalidate()
		return
	}

	if l.tex == nil {
		return
	}

	tex := args.Textures.Fetch(l.tex)
	if tex == nil {
		return
	}

	zindex := args.Position.Z
	min := args.Position.XY().Add(l.Position())
	max := min.Add(l.texSize)
	quads.Push(widget.Quad{
		Min:   min,
		Max:   max,
		MinUV: vec2.Zero,
		MaxUV: vec2.One,
		Color: [4]color.T{
			l.color,
			l.color,
			l.color,
			l.color,
		},
		ZIndex:  zindex,
		Texture: uint32(tex.ID),
	})

	// todo: cursor
}

func (l *label) measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size {
	req := flex.Size{
		Width:  l.texSize.X / l.scale,
		Height: l.texSize.Y / l.scale,
	}

	if widthMode == flex.MeasureModeExactly {
		req.Width = width
	}
	if widthMode == flex.MeasureModeAtMost {
		req.Width = math.Min(req.Width, width)
	}

	return req
}

//
// Events
//

func (l *label) FocusEvent() {
	l.state.Focused = true

}

func (l *label) BlurEvent() {
	l.state.Focused = false
	if l.props.OnBlur != nil {
		l.props.OnBlur()
	}
}

func (l *label) MouseEvent(e mouse.Event) {
	absolutePos := widget.AbsolutePosition(l.flex)
	target := e.Position().Sub(absolutePos)
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
		l.invalidate()
		l.props.OnChange(l.text)
	}
	if e.Action() == keys.Press || e.Action() == keys.Repeat {
		switch e.Code() {
		case keys.Backspace:
			str := utf8string.NewString(l.text)
			if l.cursor > 0 {
				l.cursor--
				l.text = str.Slice(0, l.cursor) + str.Slice(l.cursor+1, str.RuneCount())
				l.invalidate()
				l.props.OnChange(l.text)
			}

		case keys.Delete:
			str := utf8string.NewString(l.text)
			if l.cursor < str.RuneCount() {
				l.text = str.Slice(0, l.cursor) + str.Slice(l.cursor+1, str.RuneCount())
				l.invalidate()
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
