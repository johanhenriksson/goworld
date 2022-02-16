package deferred

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_framebuffer"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/framebuffer"

	ogl "github.com/go-gl/gl/v4.1-core/gl"
)

// GeometryPass draws the scene geometry to a G-buffer
type GeometryPass struct {
	Buffer framebuffer.Geometry
}

// NewGeometryPass sets up a geometry pass.
func NewGeometryPass() *GeometryPass {
	p := &GeometryPass{
		Buffer: gl_framebuffer.NewGeometry(1, 1),
	}
	return p
}

// DrawPass executes the geometry pass
func (p *GeometryPass) Draw(args render.Args, scene object.T) {
	p.Buffer.Bind()
	defer p.Buffer.Unbind()
	p.Buffer.Resize(args.Viewport.Width, args.Viewport.Height)

	// setup rendering
	render.Blend(false)
	render.CullFace(render.CullFront)
	render.DepthOutput(true)
	render.DepthTest(true)
	ogl.DepthFunc(ogl.GREATER)

	render.ClearWith(color.Black)
	render.ClearDepth()

	// todo: frustum culling
	// lets not draw stuff thats behind us at the very least
	// ... things need bounding boxes though.

	objects := query.New[DeferredDrawable]().Collect(scene)
	for _, drawable := range objects {
		if err := drawable.DrawDeferred(args.Apply(drawable.Object().Transform().World())); err != nil {
			fmt.Printf("deferred draw error in object %s: %s\n", drawable.Object().Name(), err)
		}
	}

	if err := gl.GetError(); err != nil {
		fmt.Println("uncaught GL error during geometry pass:", err)
	}
}
