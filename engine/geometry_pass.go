package engine

import (
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type DeferredDrawable interface {
	DrawDeferred(args DrawArgs)
}

// GeometryPass draws the scene geometry to a G-buffer
type GeometryPass struct {
	Buffer *render.GeometryBuffer
}

// Resize is called on window resize. Should update any window size-dependent buffers
func (p *GeometryPass) Resize(width, height int) {
	// recreate gbuffer
	p.Buffer.Resize(width, height)
}

// NewGeometryPass sets up a geometry pass.
func NewGeometryPass(bufferWidth, bufferHeight int) *GeometryPass {
	p := &GeometryPass{
		Buffer: render.CreateGeometryBuffer(bufferWidth, bufferHeight),
	}
	return p
}

// DrawPass executes the geometry pass
func (p *GeometryPass) Draw(scene *Scene) {
	p.Buffer.Bind()
	render.ClearWith(render.Black)
	render.ClearDepth()

	// kind-of hack to clear the diffuse buffer separately
	// allows us to clear with the camera background color
	// other buffers need to be zeroed. or???
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0) // use only diffuse buffer

	p.Buffer.DrawBuffers()

	// setup rendering
	render.Blend(false)
	render.CullFace(render.CullBack)
	render.DepthOutput(true)

	query := object.NewQuery(func(c object.Component) bool {
		_, ok := c.(DeferredDrawable)
		return ok
	})
	scene.Collect(&query)

	args := scene.Camera.DrawArgs()
	for _, component := range query.Results {
		drawable := component.(DeferredDrawable)
		drawable.DrawDeferred(args.Apply(component.Parent().Transform()))
	}

	p.Buffer.Unbind()
}
