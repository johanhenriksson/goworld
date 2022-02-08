package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/engine/deferred"
	"github.com/johanhenriksson/goworld/engine/effect"
	"github.com/johanhenriksson/goworld/render/color"
)

// PassMap maps names to Render Passes
type PassMap map[string]DrawPass

// Renderer holds references to the Scene and is responsible for executing render passes in order
type Renderer struct {
	Passes []DrawPass

	Pre      *PrePass
	Output   *OutputPass
	Geometry *deferred.GeometryPass
	Light    *deferred.LightPass
	Forward  *ForwardPass
	SSAO     *effect.SSAOPass
	Colors   *effect.ColorPass
	Lines    *LinePass

	passMap PassMap
	window  window.T
}

// NewRenderer instantiates a new rendering pipeline.
func NewRenderer(window window.T) *Renderer {
	r := &Renderer{
		passMap: make(PassMap),
		window:  window,
	}

	// deferred rendering pass
	r.Geometry = deferred.NewGeometryPass()
	r.Light = deferred.NewLightPass(r.Geometry.Buffer)

	// forward pass
	r.Forward = NewForwardPass(r.Geometry.Buffer, r.Light.Output)

	// postprocess and output
	r.SSAO = effect.NewSSAOPass(r.Geometry.Buffer, effect.SSAOSettings{
		Samples: 16,
		Radius:  0.1,
		Bias:    0.03,
		Power:   2.0,
		Scale:   2,
	})

	white := assets.GetColorTexture(color.White)
	white.Height()

	// r.Colors = effect.NewColorPass(r.Light.Output, "saturated", white)
	r.Colors = effect.NewColorPass(r.Light.Output, "saturated", r.SSAO.Gaussian.Output)
	r.Output = NewOutputPass(r.Colors.Output.Texture(), r.Geometry.Buffer.Depth())

	// lines
	r.Lines = NewLinePass()

	r.Passes = []DrawPass{
		r.Pre,
		r.Geometry,
		r.Light,
		r.Forward,
		r.SSAO,
		r.Colors,
		r.Output,
		r.Lines,
	}

	return r
}

// Append a new render pass
func (r *Renderer) Append(name string, pass DrawPass) {
	if len(name) == 0 {
		panic(fmt.Errorf("render passes must have names"))
	}

	r.Passes = append(r.Passes, pass)
	if len(name) > 0 {
		r.passMap[name] = pass
	}
}

// Get render pass by name
func (r *Renderer) Get(name string) DrawPass {
	return r.passMap[name]
}

// Draw the world.
func (r *Renderer) Draw(scene object.T) {
	camera := query.New[camera.T]().First(scene)

	args := CreateRenderArgs(r.window, camera)
	for _, pass := range r.Passes {
		pass.Draw(args, scene)
	}
}
