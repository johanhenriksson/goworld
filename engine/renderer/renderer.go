package renderer

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type T interface {
	Draw(args render.Args, scene object.T)
	Buffers() pass.BufferOutput
	Destroy()
}

type vkrenderer struct {
	Pre      *pass.PrePass
	Shadows  pass.ShadowPass
	Geometry *pass.GeometryPass
	Forward  *pass.ForwardPass
	Output   *pass.OutputPass
	Lines    *pass.LinePass
	GUI      *pass.GuiPass

	target vulkan.Target
}

func New(target vulkan.Target, geometryPasses, shadowPasses []pass.DeferredSubpass) T {
	r := &vkrenderer{
		target: target,
	}

	r.Pre = &pass.PrePass{}
	r.Shadows = pass.NewShadowPass(target, shadowPasses)
	r.Geometry = pass.NewGeometryPass(target, r.Shadows, geometryPasses)
	r.Forward = pass.NewForwardPass(target, r.Geometry.GeometryBuffer, r.Geometry.Completed())
	r.Output = pass.NewOutputPass(target, r.Geometry, r.Forward.Completed())
	r.Lines = pass.NewLinePass(target, r.Output, r.Geometry, r.Output.Completed())
	r.GUI = pass.NewGuiPass(target, r.Lines)

	return r
}

func (r *vkrenderer) Draw(args render.Args, scene object.T) {
	// render passes can be partially parallelized by dividing them into two parts,
	// recording and submission. queue submits must happen in order, so that semaphores
	// behave as expected. however, the actual recording of the command buffer can run
	// concurrently.
	//
	// to allow this, MeshCache and TextureCache must also be made thread safe, since
	// they currently work in a blocking manner.

	r.Pre.Draw(args, scene)
	r.Shadows.Draw(args, scene)
	r.Geometry.Draw(args, scene)
	r.Forward.Draw(args, scene)
	r.Output.Draw(args, scene)
	r.Lines.Draw(args, scene)
	r.GUI.Draw(args, scene)
}

func (r *vkrenderer) Buffers() pass.BufferOutput {
	return r.Geometry.GeometryBuffer
}

func (r *vkrenderer) Destroy() {
	r.target.Device().WaitIdle()

	r.GUI.Destroy()
	r.Lines.Destroy()
	r.Shadows.Destroy()
	r.Geometry.Destroy()
	r.Forward.Destroy()
	r.Output.Destroy()
}
