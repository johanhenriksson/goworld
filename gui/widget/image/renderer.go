package image

import (
	"github.com/johanhenriksson/goworld/gui/quad"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Renderer interface {
	widget.Renderer[T]

	SetSize(vec2.T)
	SetImage(texture.Ref)
	SetInvert(bool)
	SetColor(color.T)
}

type renderer struct {
	tint   color.T
	invert bool
	tex    texture.Ref

	invalid bool
	size    vec2.T
	quad    quad.T
	uvs     quad.UV
}

func NewRenderer(key string) Renderer {
	return &renderer{
		tint:    color.White,
		tex:     nil,
		uvs:     quad.DefaultUVs,
		invalid: true,
		quad:    quad.New(key, quad.Props{}),
	}
}

func (r *renderer) SetSize(size vec2.T) {
	r.invalid = r.invalid || size != r.size
	r.size = size
}

func (r *renderer) SetImage(tex texture.Ref) {
	if tex == nil {
		tex = nil
	}
	r.invalid = r.invalid || tex != r.tex
	r.tex = tex
}

func (r *renderer) SetColor(tint color.T) {
	if tint == color.None {
		tint = color.White
	}
	r.invalid = r.invalid || tint != r.tint
	r.tint = tint
}

func (r *renderer) SetInvert(invert bool) {
	r.invalid = r.invalid || invert != r.invert
	r.invert = invert
}

func (r *renderer) Draw(args widget.DrawArgs, image T) {
	if r.tex == nil {
		// nothing to render
		return
	}

	r.SetSize(image.Size())

	if r.invalid {
		uvs := r.uvs
		if r.invert {
			uvs = uvs.Inverted()
		}

		r.quad.Update(quad.Props{
			UVs:   uvs,
			Size:  r.size,
			Color: r.tint,
		})
		r.invalid = false
	}

	// fetch resources
	tex := args.Textures.Fetch(r.tex)
	if tex == 0 {
		return
	}
	mesh := args.Meshes.Fetch(r.quad.Mesh())
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
