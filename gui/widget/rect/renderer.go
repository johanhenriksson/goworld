package rect

import (
	"github.com/johanhenriksson/goworld/gui/quad"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Renderer interface {
	widget.Renderer[T]

	SetColor(color.T)
}

type renderer struct {
	mesh    quad.T
	size    vec2.T
	color   color.T
	uvs     quad.UV
	invalid bool
}

func NewRenderer(key string) Renderer {
	return &renderer{
		uvs:     quad.DefaultUVs,
		mesh:    quad.New(key, quad.Props{}),
		invalid: true,
	}
}

func (r *renderer) SetSize(size vec2.T) {
	r.invalid = r.invalid || size != r.size
	r.size = size
}

func (r *renderer) SetColor(clr color.T) {
	r.invalid = r.invalid || clr != r.color
	r.color = clr
}

func (r *renderer) drawSelf(args widget.DrawArgs, rect T) {
	// dont draw anything if its transparent anyway
	if r.color.A <= 0 {
		return
	}

	r.SetSize(rect.Size())

	if r.invalid {
		r.mesh.Update(quad.Props{
			UVs:   r.uvs,
			Size:  r.size,
			Color: r.color,
		})
		r.invalid = false
	}

	mesh := args.Meshes.Fetch(r.mesh.Mesh())
	if mesh == nil {
		// if the mesh is not available, dont draw anything this frame.
		return
	}

	tex := args.Textures.Fetch(color.White)
	if tex == nil {
		return
	}

	args.Commands.Record(func(cmd command.Buffer) {
		cmd.CmdPushConstant(core1_0.StageAll, 0, &widget.Constants{
			Viewport: args.ViewProj,
			Model:    args.Transform,
			Texture:  tex.ID,
		})
		mesh.Draw(cmd, 0)
	})
}

func (r *renderer) Draw(args widget.DrawArgs, rect T) {
	r.drawSelf(args, rect)

	for _, child := range rect.Children() {
		// calculate child transform
		// try to fix the position to an actual pixel
		// pos := vec3.Extend(child.Position().Scaled(args.Viewport.Scale).Floor().Scaled(1/args.Viewport.Scale), -1)
		z := child.ZOffset()
		pos := vec3.Extend(child.Position(), args.Position.Z-float32(1+z))
		transform := mat4.Translate(pos)
		childArgs := args
		childArgs.Transform = transform // .Mul(&args.Transform)
		childArgs.Position = pos

		// draw child
		child.Draw(childArgs)
	}
}
