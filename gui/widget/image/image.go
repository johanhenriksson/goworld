package image

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
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
	key    string
	flex   *flex.Node
	props  Props
	handle *cache.SamplerHandle
}

func New(key string, props Props) node.T {
	return node.Builtin(key, props, nil, new)
}

func new(key string, props Props) widget.T {
	node := flex.NewNodeWithConfig(flex.NewConfig())
	node.Context = key
	img := &image{
		key:  key,
		flex: node,
	}
	img.flex.SetMeasureFunc(img.measure)
	img.Update(props)
	return img
}

//
// Widget implementation
//

func (i *image) Key() string          { return i.key }
func (i *image) Flex() *flex.Node     { return i.flex }
func (i *image) Children() []widget.T { return nil }
func (i *image) Position() vec2.T     { return vec2.New(i.flex.LayoutGetLeft(), i.flex.LayoutGetTop()) }
func (i *image) Size() vec2.T         { return vec2.New(i.flex.LayoutGetWidth(), i.flex.LayoutGetHeight()) }

func (i *image) SetChildren(children []widget.T) {
	if len(children) > 0 {
		panic("images may not have children")
	}
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
// Styles
//

func (i *image) Image() texture.Ref { return i.props.Image }

//
// Draw
//

func (i *image) Draw(args widget.DrawArgs, quads *widget.QuadBuffer) {
	tex, texExists := args.Textures.TryFetch(i.props.Image)
	if !texExists {
		return
	}

	if tex != i.handle {
		// image handle changed, redo layout
		i.handle = tex
		i.Flex().MarkDirty()
	}

	zindex := args.Position.Z
	min := args.Position.XY().Add(i.Position())
	max := min.Add(i.Size())

	quads.Push(widget.Quad{
		Min:   min,
		Max:   max,
		MinUV: vec2.Zero,
		MaxUV: vec2.One,
		Color: [4]color.T{ // todo: add tint prop
			color.White,
			color.White,
			color.White,
			color.White,
		},
		ZIndex:  zindex,
		Texture: uint32(tex.ID),
	})
}

//
// Events
//

func (i *image) Destroy() {
}

func (l *image) MouseEvent(e mouse.Event) {
	if e.Action() == mouse.Press && l.props.OnClick != nil {
		l.props.OnClick(e)
		e.Consume()
	}
}
