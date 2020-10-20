package engine

import (
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type DeferredDrawable interface {
	DrawDeferred(args DrawArgs)
}

// GeometryPass draws the scene geometry to a G-buffer
type GeometryPass struct {
	Buffer *render.GeometryBuffer
	queue  *DrawQueue
}

func (p *GeometryPass) Type() render.Pass {
	return render.Geometry
}

// NewGeometryPass sets up a geometry pass.
func NewGeometryPass(bufferWidth, bufferHeight int) *GeometryPass {
	p := &GeometryPass{
		Buffer: render.CreateGeometryBuffer(bufferWidth, bufferHeight),
		queue:  NewDrawQueue(),
	}
	return p
}

// DrawPass executes the geometry pass
func (p *GeometryPass) DrawPass(scene *Scene) {
	p.Buffer.Bind()
	render.Clear()
	render.ClearDepth()

	// kind-of hack to clear the diffuse buffer separately
	// allows us to clear with the camera background color
	// other buffers need to be zeroed. or???
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0) // use only diffuse buffer
	render.ClearWith(scene.Camera.Clear)

	p.Buffer.DrawBuffers()

	// setup rendering
	render.Blend(false)
	render.CullFace(render.CullBack)
	render.DepthOutput(true)

	p.queue.Clear()
	scene.Collect(p)

	for _, cmd := range p.queue.items {
		drawable := cmd.Component.(DeferredDrawable)
		drawable.DrawDeferred(cmd.Args)
	}

	p.Buffer.Unbind()
}

func (p *GeometryPass) Visible(c Component, args DrawArgs) bool {
	_, ok := c.(DeferredDrawable)
	return ok
}

func (p *GeometryPass) Queue(c Component, args DrawArgs) {
	p.queue.Add(c, args)
}
