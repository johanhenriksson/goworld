package rect

import (
	"log"

	"github.com/johanhenriksson/goworld/gui/quad"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/texture"

	vk "github.com/vulkan-go/vulkan"
)

type Renderer interface {
	widget.Renderer[T]

	SetColor(color.T)
}

type renderer struct {
	tex     texture.T
	mesh    quad.T
	size    vec2.T
	color   color.T
	uvs     quad.UV
	invalid bool
}

func NewRenderer() Renderer {
	return &renderer{
		uvs:     quad.DefaultUVs,
		mesh:    quad.New(quad.Props{}),
		invalid: true,
	}
}

func (r *renderer) SetSize(size vec2.T) {
	// log.Println("rect size update", size != r.size, size, r.size)
	r.invalid = r.invalid || size != r.size
	r.size = size
}

func (r *renderer) SetColor(clr color.T) {
	log.Println("rect color update", clr != r.color, clr, r.color)
	r.invalid = r.invalid || clr != r.color
	r.color = clr
}

func (r *renderer) Draw(args widget.DrawArgs, rect T) {

	// set correct blending
	// render.BlendMultiply()

	// render.Scissor(frame.Position(), frame.Size())

	// dont draw anything if its transparent anyway
	if r.color.A > 0 {
		r.SetSize(rect.Size())

		if r.invalid {
			log.Println("updating rect", rect.Key())
			r.mesh.Update(quad.Props{
				UVs:   r.uvs,
				Size:  r.size,
				Color: r.color,
			})
			r.invalid = false
		}

		mesh := args.Meshes.Fetch(r.mesh.Mesh())
		if mesh != nil {
			args.Commands.Record(func(cmd command.Buffer) {
				cmd.CmdPushConstant(vk.ShaderStageAll, 0, &widget.Constants{
					Viewport: args.ViewProj,
					Model:    args.Transform,
					Texture:  0,
				})
				mesh.Draw(cmd, 0)
			})
		}
	}

	for _, child := range rect.Children() {
		// calculate child tranasform
		// try to fix the position to an actual pixel
		// pos := vec3.Extend(child.Position().Scaled(args.Viewport.Scale).Floor().Scaled(1/args.Viewport.Scale), -1)
		pos := vec3.Extend(child.Position(), args.Position.Z-1)
		transform := mat4.Translate(pos)
		childArgs := args
		childArgs.Transform = transform // .Mul(&args.Transform)
		childArgs.Position = pos

		// draw child
		child.Draw(childArgs)
	}

	// render.ScissorDisable()
}

func (r *renderer) Destroy() {
	//  todo: clean up mesh, texture
}
