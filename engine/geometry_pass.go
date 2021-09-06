package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/render"
)

type DeferredDrawable interface {
	DrawDeferred(render.Args)
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
func (p *GeometryPass) Draw(scene scene.T) {
	p.Buffer.Bind()
	render.ClearWith(render.Black)
	render.ClearDepth()

	// setup rendering
	render.Blend(false)
	render.CullFace(render.CullBack)
	render.DepthOutput(true)

	query := object.NewQuery(DeferredDrawableQuery)
	scene.Collect(&query)

	args := ArgsFromCamera(scene.Camera())
	for _, component := range query.Results {
		drawable := component.(DeferredDrawable)
		drawable.DrawDeferred(args.Apply(component.Object().Transform().World()))
	}

	p.Buffer.Unbind()
}

// DeferedDrawableQuery is an object query predicate that matches any component
// that implements the DeferredDrawable interface.
func DeferredDrawableQuery(c object.Component) bool {
	_, ok := c.(DeferredDrawable)
	return ok
}
