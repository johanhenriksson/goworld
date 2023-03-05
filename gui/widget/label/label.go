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

	"github.com/kjk/flex"
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
	key   string
	flex  *flex.Node
	props Props
	scale float32
	state style.State
	text  *Text

	fontSize   int
	fontName   string
	font       font.T
	color      color.T
	highlight  color.T
	lineHeight float32
	size       vec2.T
}

func New(key string, props Props) node.T {
	return node.Builtin(key, props, nil, new)
}

func new(key string, props Props) T {
	node := flex.NewNodeWithConfig(flex.NewConfig())
	node.Context = key
	text := NewText(props.Text)

	lbl := &label{
		key:   key,
		flex:  node,
		scale: 1,
		text:  text,

		highlight:  color.RGBA(0, 0, 0, 0.3),
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

	// refresh font if required
	if l.font == nil {
		fontName := l.fontName
		if fontName == "" {
			fontName = DefaultFont.Name
		}
		fontSize := l.fontSize
		if fontSize == 0 {
			fontSize = DefaultFont.Size
		}
		l.font = assets.GetFont(fontName, fontSize, l.scale)
	}

	// recalculate text size
	l.size = l.font.Measure(l.text.String(), font.Args{LineHeight: l.lineHeight})
}

func (l *label) editable() bool { return l.props.OnChange != nil }

func (l *label) setText(text string) {
	if text == l.text.String() {
		return
	}

	l.text.SetText(text)
	l.invalidate()
}

//
// Styles
//

func (l *label) SetFont(font style.Font) {
	if font.Name == l.fontName && font.Size == l.fontSize {
		return
	}
	l.fontName = font.Name
	l.fontSize = font.Size
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
func (l *label) Cursor() int  { return l.text.cursor }

//
// Draw
//

func (l *label) Draw(args widget.DrawArgs, quads *widget.QuadBuffer) {
	l.text.UpdateBlink(args.Delta)

	if l.props.Style.Hidden {
		return
	}

	if args.Viewport.Scale != l.scale {
		// ui scale has changed
		l.scale = args.Viewport.Scale
		l.font = nil
		l.invalidate()
		return
	}

	if l.font == nil {
		// no font yet
		return
	}

	zindex := args.Position.Z
	origin := args.Position.XY().Add(l.Position())

	// render glyph quads
	pos := origin.Add(vec2.New(0, 0.8*l.font.Size()))
	for _, r := range l.text.String() {
		glyph, err := l.font.Glyph(r)
		if err != nil {
			panic(err)
		}
		if glyph.Size.Y == 0 {
			// whitespace
			pos.X += glyph.Advance
			continue
		}

		handle := args.Textures.Fetch(glyph)
		if handle != nil {
			min := pos.Add(glyph.Bearing).Floor()
			max := min.Add(glyph.Size)
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
				Texture: uint32(handle.ID),
			})
		}

		pos.X += glyph.Advance
	}

	// selection
	if l.text.HasSelection() {
		// measure selection size & position
		// todo: could be cached, perhaps not necessary
		args := font.Args{LineHeight: l.lineHeight}
		startIdx, _ := l.text.SelectedRange()
		start := l.font.Measure(l.text.Slice(0, startIdx), args)
		selectSize := l.font.Measure(l.text.Selection(), args)
		length := math.Max(selectSize.X, 1)

		// highlight quad
		min := origin.Add(vec2.New(start.X, 0))
		max := min.Add(vec2.New(length, l.lineHeight*l.font.Size()))
		quads.Push(widget.Quad{
			Min:     min,
			Max:     max,
			Color:   [4]color.T{l.highlight, l.highlight, l.highlight, l.highlight},
			ZIndex:  zindex + 0.1,
			Texture: 0,
		})
	}

	// cursor
	if l.state.Focused && !l.text.HasSelection() && l.text.Blink() {
		// measure cursor position
		args := font.Args{LineHeight: l.lineHeight}
		cursorIdx, _ := l.text.SelectedRange()
		cursorPos := l.font.Measure(l.text.Slice(0, cursorIdx), args)

		// cursor quad
		min := origin.Add(vec2.New(cursorPos.X, 0))
		max := min.Add(vec2.New(1, l.lineHeight*l.font.Size()))
		quads.Push(widget.Quad{
			Min:     min,
			Max:     max,
			Color:   [4]color.T{color.Black, color.Black, color.Black, color.Black},
			ZIndex:  zindex + 0.1,
			Texture: 0,
		})
	}
}

func (l *label) measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size {
	req := flex.Size{
		Width:  l.size.X,
		Height: l.size.Y,
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

func (l *label) Destroy() {
	if l.state.Focused {
		keys.Focus(nil)
	}
}

func (l *label) FocusEvent() {
	l.state.Focused = true
	// todo: move cursor to mouse position?
}

func (l *label) BlurEvent() {
	l.text.Deselect()
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
		l.text.Insert(string(e.Character()))
		l.invalidate()
		l.props.OnChange(l.text.String())
	}
	if e.Action() == keys.Press || e.Action() == keys.Repeat {
		switch e.Code() {
		case keys.Backspace:
			if l.text.DeleteBackward() {
				l.invalidate()
				l.props.OnChange(l.text.String())
			}

		case keys.Delete:
			if l.text.DeleteForward() {
				l.invalidate()
				l.props.OnChange(l.text.String())
			}

		case keys.LeftArrow:
			if e.Modifier(keys.Shift) {
				l.text.SelectLeft()
			} else {
				l.text.CursorLeft()
			}

		case keys.RightArrow:
			if e.Modifier(keys.Shift) {
				l.text.SelectRight()
			} else {
				l.text.CursorRight()
			}

		case keys.U:
			// ctrl+u clears text
			if e.Modifier(keys.Ctrl) {
				if l.text.Clear() {
					l.props.OnChange("")
				}
			}
		}
	}
}
