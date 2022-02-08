package image

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/kjk/flex"
)

type T interface {
	widget.T
	Tint() color.T
	Image() texture.T
}

type Props struct {
	Style   style.Sheet
	Image   texture.T
	Tint    color.T
	Invert  bool
	OnClick mouse.Callback
}

type image struct {
	widget.T
	props    Props
	renderer Renderer
}

func New(key string, props Props) node.T {
	return node.Builtin(key, props, nil, new)
}

func new(key string, props Props) T {
	img := &image{
		T:        widget.New(key),
		renderer: &renderer{},
	}
	img.Update(props)
	return img
}

func (i *image) Props() any { return i.props }

func (i *image) Update(props any) {
	new := props.(Props)

	// default
	if new.Tint == color.None {
		new.Tint = color.White
	}
	if new.Image == nil {
		new.Image = assets.DefaultTexture()
	}

	imageChanged := new.Image != i.props.Image
	invalidated := imageChanged

	styleChanged := new.Style != i.props.Style

	// update props
	i.props = new

	if styleChanged {
		i.SetStyle(new.Style)
	}

	if invalidated {
		i.Flex().MarkDirty()
	}
}

// prop accessors

func (i *image) Image() texture.T { return i.props.Image }
func (i *image) Tint() color.T    { return i.props.Tint }

func (i *image) Draw(args render.Args) {
	i.T.Draw(args)
	i.renderer.Draw(args, i, &i.props)
}

func (i *image) Flex() *flex.Node {
	node := i.T.Flex()
	node.SetMeasureFunc(i.measure)
	return node
}

func (i *image) measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size {
	// todo: consider constraints
	size := i.props.Image.Size()
	aspect := size.X / size.Y
	return flex.Size{
		Width:  width,
		Height: width / aspect,
	}
}

//
// Events
//

func (l *image) MouseEvent(e mouse.Event) {
	if e.Action() == mouse.Press && l.props.OnClick != nil {
		l.props.OnClick(e)
		e.Consume()
	}
}
