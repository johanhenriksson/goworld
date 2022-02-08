package label

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/font"
	"github.com/kjk/flex"
)

type T interface {
	widget.T
	Font() font.T
	Text() string
	LineHeight() float32
}

type Props struct {
	Style      style.Sheet
	Text       string
	Font       font.T
	Size       int
	LineHeight float32
	OnClick    mouse.Callback
}

type label struct {
	widget.T
	props    Props
	renderer Renderer
	size     vec2.T
}

func New(key string, props Props) node.T {
	if props.Size == 0 {
		props.Size = 12
	}
	if props.LineHeight == 0 {
		props.LineHeight = 0
	}
	if props.Font == nil {
		props.Font = assets.DefaultFont()
	}
	return node.Builtin(key, props, nil, new)
}

func new(key string, props Props) T {
	lbl := &label{
		T:        widget.New(key),
		renderer: &renderer{},
	}
	lbl.Update(props)
	return lbl
}

func (l *label) Size() vec2.T { return l.T.Size() }

func (l *label) Props() any { return l.props }
func (l *label) Update(props any) {
	new := props.(Props)

	textChanged := new.Text != l.props.Text
	sizeChanged := new.Size != l.props.Size
	heightChanged := new.LineHeight != l.props.LineHeight
	styleChanged := new.Style != l.props.Style
	invalidated := textChanged || sizeChanged || heightChanged

	l.props = new

	if styleChanged {
		l.SetStyle(new.Style)
	}

	if invalidated || styleChanged {
		l.Flex().MarkDirty()
	}
}

// prop accessors

func (l *label) Font() font.T        { return l.props.Font }
func (l *label) Text() string        { return l.props.Text }
func (l *label) LineHeight() float32 { return l.props.LineHeight }

func (l *label) Draw(args render.Args) {
	l.T.Draw(args)
	l.renderer.Draw(args, l, &l.props)
}

func (l *label) Flex() *flex.Node {
	node := l.T.Flex()
	node.SetMeasureFunc(l.measure)
	return node
}

func (l *label) measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size {
	// todo: consider constraints
	args := font.Args{
		LineHeight: l.props.LineHeight,
		Color:      color.White,
	}
	size := l.props.Font.Measure(l.props.Text, args).Scaled(0.5)

	return flex.Size{
		Width:  size.X,
		Height: size.Y,
	}
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
