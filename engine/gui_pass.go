package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type GuiDrawable interface {
	object.Component
	DrawUI(args render.Args, scene object.T)
}

type GuiPass struct {
}

func NewGuiPass() *GuiPass {
	return &GuiPass{}
}

func (g *GuiPass) Draw(args render.Args, scene object.T) {
	width, height := float32(args.Viewport.FrameWidth), float32(args.Viewport.FrameHeight)
	scale := width / float32(args.Viewport.Width)

	// setup viewport
	proj := mat4.OrthographicVK(0, width, height, 0, 1000, -1000)
	view := mat4.Scale(vec3.New(scale, scale, 1))
	vp := proj.Mul(&view)

	gl.DepthFunc(gl.GREATER)

	uiArgs := render.Args{
		Projection: proj,
		View:       view,
		VP:         vp,
		MVP:        vp,
		Transform:  mat4.Ident(),
		Viewport:   args.Viewport,
	}

	// query scene for gui managers
	guis := query.New[GuiDrawable]().Collect(scene)
	for _, gui := range guis {
		gui.DrawUI(uiArgs, scene)
	}
}
