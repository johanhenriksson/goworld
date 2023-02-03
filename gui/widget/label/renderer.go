package label

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/gui/quad"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/font"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/util"

	"github.com/kjk/flex"
	"github.com/vkngwrapper/core/v2/core1_0"
)

var DefaultFont = "fonts/SourceCodeProRegular.ttf"
var DefaultSize = 12
var DefaultColor = color.White
var DefaultLineHeight = float32(1.0)

type Renderer interface {
	widget.Renderer[T]

	SetText(string)
	SetFont(string)
	SetFontSize(int)
	SetFontColor(color.T)
	SetLineHeight(float32)

	Measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size
}

type renderer struct {
	key        string
	text       string
	version    int
	size       int
	fontName   string
	font       font.T
	color      color.T
	lineHeight float32

	invalidFont    bool
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
		key:            fmt.Sprintf("label:%s", util.NewUUID(8)),
		size:           DefaultSize,
		fontName:       DefaultFont,
		color:          DefaultColor,
		lineHeight:     DefaultLineHeight,
		invalidFont:    true,
		invalidTexture: true,
		invalidMesh:    true,
		scale:          2,
		uvs:            quad.DefaultUVs,
		mesh:           quad.New(quad.Props{}),
	}
}

func (r *renderer) SetText(text string) {
	r.invalidTexture = r.invalidTexture || text != r.text
	r.invalidMesh = r.invalidMesh || text != r.text
	r.text = text
}

func (r *renderer) SetFont(name string) {
	r.invalidFont = name != r.fontName
	r.invalidTexture = r.invalidTexture || r.invalidFont
	r.invalidMesh = r.invalidMesh || r.invalidFont
	r.fontName = name
}

func (r *renderer) SetFontSize(size int) {
	if size <= 0 {
		size = 12
	}
	r.invalidTexture = r.invalidTexture || size != r.size
	r.invalidMesh = r.invalidMesh || size != r.size
	r.invalidFont = r.invalidFont || size != r.size
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

	if r.invalidFont {
		r.font = assets.GetFont(r.fontName, r.size, r.scale)
		r.invalidFont = false
	}

	if r.invalidTexture {
		// (re)create label texture
		fargs := font.Args{
			LineHeight: r.lineHeight,
			Color:      color.White,
		}
		r.bounds = r.font.Measure(r.text, fargs)

		r.version++
		// img := r.font.Render(r.text, fargs)
		// r.tex = texture.ImageRef(r.key, r.version, img)
		key := fmt.Sprint(r.version, r.fontName, r.size, r.text, fargs)
		r.tex = font.Ref(key, r.version, r.font, r.text, fargs)
		log.Println("invalid:", r.text)

		r.invalidTexture = false
	}

	if r.invalidMesh {
		r.mesh.Update(quad.Props{
			Size:  r.bounds.Scaled(1 / r.scale),
			UVs:   r.uvs,
			Color: r.color,
		})
		r.invalidMesh = false
	}

	// resize mesh if needed
	// if !label.Size().ApproxEqual(r.size) {
	// 	fmt.Println("label size", label.Size())
	// 	r.mesh.SetSize(label.Size())
	// 	r.size = label.Size()
	// }

	// can the we use the gl viewport to clip anything out of bounds?

	// we can center the label on the mesh by modifying the uvs
	// scale := label.Size().Div(r.bounds)

	tex := args.Textures.Fetch(r.tex)
	if tex == 0 {
		return
	}
	mesh := args.Meshes.Fetch(r.mesh.Mesh())
	if mesh == nil {
		return
	}

	args.Commands.Record(func(cmd command.Buffer) {
		cmd.CmdPushConstant(core1_0.StageAll, 0, &widget.Constants{
			Viewport: args.ViewProj,
			Model:    args.Transform,
			Texture:  tex,
		})
		mesh.Draw(cmd, 0)
	})
}

func (r *renderer) Measure(node *flex.Node, width float32, widthMode flex.MeasureMode, height float32, heightMode flex.MeasureMode) flex.Size {
	if r.font == nil {
		return flex.Size{}
	}

	size := r.font.Measure(r.text, font.Args{
		LineHeight: r.lineHeight,
		Color:      color.White,
	})

	// size = size.Scaled(1 / r.scale)

	return flex.Size{
		Width:  size.X / r.scale,
		Height: size.Y / r.scale,
	}
}
