package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_framebuffer"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/framebuffer"
)

type DeferredDrawable interface {
	DrawDeferred(render.Args)
}

// GeometryPass draws the scene geometry to a G-buffer
type GeometryPass struct {
	Buffer framebuffer.Geometry
}

// NewGeometryPass sets up a geometry pass.
func NewGeometryPass(bufferWidth, bufferHeight int) *GeometryPass {
	p := &GeometryPass{
		Buffer: gl_framebuffer.NewGeometry(bufferWidth, bufferHeight),
	}
	return p
}

// DrawPass executes the geometry pass
func (p *GeometryPass) Draw(args render.Args, scene scene.T) {
	p.Buffer.Bind()
	defer p.Buffer.Unbind()
	p.Buffer.Resize(args.Viewport.FrameWidth, args.Viewport.FrameHeight)

	// setup rendering
	render.Blend(false)
	render.CullFace(render.CullBack)
	render.DepthOutput(true)
	render.DepthTest(true)

	render.ClearWith(color.Black)
	render.ClearDepth()

	query := object.NewQuery(DeferredDrawableQuery)
	scene.Collect(&query)

	args = ArgsWithCamera(args, scene.Camera())
	for _, component := range query.Results {
		drawable := component.(DeferredDrawable)
		drawable.DrawDeferred(args.Apply(component.Object().Transform().World()))
	}
}

// DeferedDrawableQuery is an object query predicate that matches any component
// that implements the DeferredDrawable interface.
func DeferredDrawableQuery(c object.Component) bool {
	_, ok := c.(DeferredDrawable)
	return ok
}
