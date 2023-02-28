package image

import (
	"log"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/color"
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
	OnClick mouse.Callback
}

type image struct {
	widget.T
	props  Props
	handle *cache.SamplerHandle
}

func New(key string, props Props) node.T {
	return node.Builtin(key, props, nil, new)
}

func new(w widget.T, props Props) T {
	img := &image{
		T: w,
	}
	w.Flex().SetMeasureFunc(img.measure)
	img.Update(props)
	return img
}

func (i *image) Props() any { return i.props }

func (i *image) Update(props any) {
	new := props.(Props)
	styleChanged := new.Style != i.props.Style
	i.props = new

	if styleChanged {
		new.Style.Apply(i, style.State{})
		i.Flex().MarkDirty()
	}
}

// prop accessors

func (i *image) Image() texture.Ref { return i.props.Image }

func (i *image) Draw(args widget.DrawArgs, quads *widget.QuadBuffer) {
	tex := args.Textures.Fetch(i.props.Image)
	if tex != nil {
		if tex != i.handle {
			// image handle changed, redo layout
			i.Flex().MarkDirty()
			i.handle = tex
		}

		log.Println("draw image with texture", tex.ID, "and size", i.Size())
		quads.Push(widget.Quad{
			Min:     args.Position.XY(),
			Max:     args.Position.XY().Add(i.Size()),
			MinUV:   vec2.Zero,
			MaxUV:   vec2.One,
			Color:   color.White, // todo: add tint prop
			ZIndex:  args.Position.Z,
			Texture: uint32(tex.ID),
		})
	}
}

func (i *image) measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size {
	if i.handle == nil {
		return flex.Size{}
	}

	// todo: consider constraints
	w := float32(i.handle.Texture.Image().Width())
	h := float32(i.handle.Texture.Image().Height())
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
