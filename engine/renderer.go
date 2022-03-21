package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/deferred"
	"github.com/johanhenriksson/goworld/engine/effect"
	"github.com/johanhenriksson/goworld/render"
)

type Renderer interface {
	Draw(args render.Args, scene object.T)
	Buffers() BufferOutput
	Destroy()
}

// GLRenderer holds references to the Scene and is responsible for executing render passes in order
type GLRenderer struct {
	Passes []DrawPass

	Pre      *PrePass
	Output   *OutputPass
	Geometry *deferred.GeometryPass
	Light    *deferred.LightPass
	Forward  *ForwardPass
	SSAO     *effect.SSAOPass
	Colors   *effect.ColorPass
	Lines    *LinePass
	Gui      *GuiPass

	meshes cache.Meshes
}

// NewRenderer instantiates a new rendering pipeline.
func NewRenderer() Renderer {
	r := &GLRenderer{
		meshes: cache.NewMeshes(),
	}

	// deferred rendering pass
	r.Geometry = deferred.NewGeometryPass(r.meshes)
	r.Light = deferred.NewLightPass(r.Geometry.Buffer, r.meshes)

	// forward pass
	r.Forward = NewForwardPass(r.Geometry.Buffer, r.Light.Output, r.meshes)

	// ssao pass
	r.SSAO = effect.NewSSAOPass(r.Geometry.Buffer, effect.SSAOSettings{
		Samples: 16,
		Radius:  0.3,
		Bias:    0.03,
		Power:   2.0,
		Scale:   2,
	})

	// color correction & ssao merge
	r.Colors = effect.NewColorPass(r.Light.Output, "saturated", r.SSAO.Gaussian.Output)

	// output world image
	r.Output = NewOutputPass(r.Colors.Output.Texture(), r.Geometry.Buffer.Depth())

	// lines
	r.Lines = NewLinePass(r.meshes)

	// gui
	r.Gui = NewGuiPass()

	r.Passes = []DrawPass{
		r.Pre,
		r.Geometry,
		r.Light,
		r.Forward,
		r.SSAO,
		r.Colors,
		r.Output,
		r.Lines,
		r.Gui,
	}

	return r
}

// Draw the world.
func (r *GLRenderer) Draw(args render.Args, scene object.T) {
	for _, pass := range r.Passes {
		pass.Draw(args, scene)
	}

	// reclaim mesh memory
	r.meshes.Tick()
	r.meshes.Evict()
}

func (r *GLRenderer) Buffers() BufferOutput {
	return nil
}

func (r *GLRenderer) Destroy() {

}
