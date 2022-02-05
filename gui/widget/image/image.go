package image

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

type T interface {
	widget.T
	Tint() color.T
	Image() texture.T
}

type Props struct {
	Image   texture.T
	Tint    color.T
	Invert  bool
	OnClick mouse.Callback
}

type image struct {
	widget.T
	props    *Props
	renderer Renderer
	size     vec2.T
}

func New(key string, props *Props) node.T {
	return node.Builtin(key, props, nil, new)
}

func new(key string, props *Props) T {
	img := &image{
		T:        widget.New(key),
		renderer: &renderer{},
	}
	img.Update(props)
	return img
}

func (i *image) Size() vec2.T { return i.T.Size() }

func (i *image) Props() any { return i.props }
func (i *image) Update(props any) {
	i.props = props.(*Props)
	if i.props.Tint == color.None {
		i.props.Tint = color.White
	}
	if i.props.Image == nil {
		i.props.Image = assets.DefaultTexture()
	}
	i.size = i.props.Image.Size()
}

// prop accessors

func (i *image) Image() texture.T { return i.props.Image }
func (i *image) Tint() color.T    { return i.props.Tint }

func (i *image) Draw(args render.Args) {
	i.T.Draw(args)
	i.renderer.Draw(args, i, i.props)
}

func (i *image) Width() dimension.T  { return dimension.Fixed(i.size.X) }
func (i *image) Height() dimension.T { return dimension.Fixed(i.size.Y) }

//
// Events
//

func (l *image) MouseEvent(e mouse.Event) {
	if e.Action() == mouse.Press && l.props.OnClick != nil {
		l.props.OnClick(e)
		e.Consume()
	}
}
