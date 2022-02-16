package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
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
	size := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	scale := args.Viewport.Scale
	size = size.Scaled(1 / scale)

	// setup viewport
	proj := mat4.OrthographicLH(0, size.X, size.Y, 0, 1000, -1000)
	view := mat4.Ident()
	vp := proj.Mul(&view)

	gl.DepthFunc(gl.GREATER)

	uiArgs := render.Args{
		Projection: proj,
		View:       view,
		VP:         vp,
		MVP:        vp,
		Transform:  mat4.Ident(),
		Viewport: render.Screen{
			Width:  int(size.X),
			Height: int(size.Y),
			Scale:  scale,
		},
	}

	// query scene for gui managers
	guis := query.New[GuiDrawable]().Collect(scene)
	for _, gui := range guis {
		gui.DrawUI(uiArgs, scene)
	}
}
