package label

import (
	"log"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/gui/quad"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/texture"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/font"

	"github.com/kjk/flex"
)

var DefaultSize = 12
var DefaultColor = color.White
var DefaultLineHeight = float32(1.0)

type Renderer interface {
	widget.Renderer[T]

	SetText(string)
	SetFont(font.T)
	SetFontSize(int)
	SetFontColor(color.T)
	SetLineHeight(float32)

	Measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size
}

type renderer struct {
	text       string
	size       int
	font       font.T
	color      color.T
	lineHeight float32

	invalidTexture bool
	invalidMesh    bool
	scale          float32
	bounds         vec2.T
	tex            texture.Ref
	mesh           quad.T
	uvs            quad.UV
}

func NewRenderer() Renderer {
	return &renderer{
		size:           DefaultSize,
		color:          DefaultColor,
		lineHeight:     DefaultLineHeight,
		font:           assets.DefaultFont(),
		invalidTexture: true,
		invalidMesh:    true,
		scale:          2,
		uvs:            quad.DefaultUVs,
		mesh:           quad.New(quad.Props{}),
	}
}

func (r *renderer) SetText(text string) {
	r.invalidTexture = r.invalidTexture || text != r.text
	r.text = text
}

func (r *renderer) SetFont(fnt font.T) {
	if fnt == nil {
		fnt = assets.DefaultFont()
	}
	r.invalidTexture = r.invalidTexture || fnt != r.font
	r.invalidMesh = r.invalidMesh || fnt != r.font
	r.font = fnt
}

func (r *renderer) SetFontSize(size int) {
	if size <= 0 {
		size = 12
	}
	r.invalidTexture = r.invalidTexture || size != r.size
	r.invalidMesh = r.invalidMesh || size != r.size
	r.size = size
}

func (r *renderer) SetFontColor(clr color.T) {
	if clr == color.None {
		clr = color.White
	}
	r.invalidMesh = r.invalidMesh || clr != r.color
	r.color = clr
}

func (r *renderer) SetLineHeight(lineHeight float32) {
	if lineHeight <= 0 {
		lineHeight = 1
	}
	r.invalidTexture = lineHeight != r.lineHeight
	r.invalidMesh = r.invalidMesh || lineHeight != r.lineHeight
	r.lineHeight = lineHeight
}

func (r *renderer) Draw(args widget.DrawArgs, label T) {
	r.scale = args.Viewport.Scale

	if r.text == "" {
		return
	}

	if r.invalidTexture {
		log.Println("update label texture")
		// (re)create label texture
		fargs := font.Args{
			LineHeight: r.lineHeight,
			Color:      color.White,
		}
		r.bounds = r.font.Measure(r.text, fargs)

		img := r.font.Render(r.text, fargs)
		r.tex = texture.ImageRef("label", img)

		r.invalidTexture = false
	}

	if r.invalidMesh {
		log.Println("update label mesh")
		r.mesh.Update(quad.Props{
			Size:  r.bounds.Scaled(1 / r.scale),
			UVs:   r.uvs,
			Color: r.color,
		})
		r.invalidMesh = false
	}

	// set correct blending
	// render.BlendMultiply()

	// resize mesh if needed
	// if !label.Size().ApproxEqual(r.size) {
	// 	fmt.Println("label size", label.Size())
	// 	r.mesh.SetSize(label.Size())
	// 	r.size = label.Size()
	// }

	// can the we use the gl viewport to clip anything out of bounds?

	// we can center the label on the mesh by modifying the uvs
	// scale := label.Size().Div(r.bounds)

	// r.mesh.Draw(args)
}

func (r *renderer) Destroy() {
	if r.mesh != nil {
		// r.mesh.Destroy()
		r.mesh = nil
	}
	if r.tex != nil {
		// todo: delete texture
		r.tex = nil
	}
}

func (r *renderer) Measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size {
	size := r.font.Measure(r.text, font.Args{
		LineHeight: r.lineHeight,
		Color:      color.White,
	})

	size = size.Scaled(1 / r.scale)

	return flex.Size{
		Width:  size.X,
		Height: size.Y,
	}
}
