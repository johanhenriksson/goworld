package image

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/kjk/flex"
)

type T interface {
	widget.T
	Image() texture.Ref
}

type Props struct {
	Style   Style
	Image   texture.Ref
	Invert  bool
	OnClick mouse.Callback
}

type image struct {
	widget.T
	Renderer
	props   Props
	texture texture.T
}

func New(key string, props Props) node.T {
	return node.Builtin(key, props, nil, new)
}

func new(w widget.T, props Props) T {
	img := &image{
		T: w,
	}
	img.Renderer = NewRenderer(w.Key(), img.onLoad)
	w.Flex().SetMeasureFunc(img.measure)
	img.Update(props)
	return img
}

func (i *image) Props() any { return i.props }

func (i *image) Update(props any) {
	new := props.(Props)
	styleChanged := new.Style != i.props.Style
	i.props = new

	i.Renderer.SetImage(new.Image)
	i.Renderer.SetInvert(new.Invert)

	if styleChanged {
		new.Style.Apply(i, style.State{})
		i.Flex().MarkDirty()
	}
}

// prop accessors

func (i *image) Image() texture.Ref { return i.props.Image }

func (i *image) Draw(args widget.DrawArgs) {
	i.T.Draw(args)
	i.Renderer.Draw(args, i)
}

func (i *image) onLoad(tex texture.T) {
	// update active texture reference & trigger layout update
	i.texture = tex
	i.Flex().MarkDirty()
}

func (i *image) measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size {
	if i.texture == nil {
		return flex.Size{}
	}

	// todo: consider constraints
	w := float32(i.texture.Image().Width())
	h := float32(i.texture.Image().Height())
	aspect := w / h
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
