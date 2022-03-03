package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_framebuffer"
	"github.com/johanhenriksson/goworld/render/framebuffer"

	ogl "github.com/go-gl/gl/v4.1-core/gl"
)

// ForwardPass holds information required to perform a forward rendering pass.
type ForwardPass struct {
	output  framebuffer.Color
	gbuffer framebuffer.Geometry
	fbo     framebuffer.T
	meshes  cache.Meshes
}

// NewForwardPass sets up a forward pass.
func NewForwardPass(gbuffer framebuffer.Geometry, output framebuffer.Color, meshes cache.Meshes) *ForwardPass {
	// the forward pass renders into the output of the final deferred pass.
	// it reuses the normal, position and depth buffers and writes new data according to what is rendered
	// this ensures that we have complete information in those buffers for later passes
	fbo := gl_framebuffer.NewGeometry(gbuffer.Width(), gbuffer.Height())
	fbo.AttachBuffer(ogl.COLOR_ATTACHMENT0, output.Texture())
	fbo.AttachBuffer(ogl.COLOR_ATTACHMENT1, gbuffer.Normal())
	fbo.AttachBuffer(ogl.COLOR_ATTACHMENT2, gbuffer.Position())
	fbo.AttachBuffer(ogl.DEPTH_ATTACHMENT, gbuffer.Depth())

	return &ForwardPass{
		fbo:     fbo,
		output:  output,
		gbuffer: gbuffer,
		meshes:  meshes,
	}
}

// DrawPass executes the forward pass
func (p *ForwardPass) Draw(args render.Args, scene object.T) {

	// setup rendering
	render.Blend(true)
	render.BlendMultiply()
	render.CullFace(render.CullFront)

	// todo: depth sorting
	// there is finally a decent way of doing it!!
	// now we just need a way to compute the distance from an object to the camera
	// ... and a way to sort the queue

	// todo: frustum culling
	// lets not draw stuff thats behind us at the very least
	// ... things need bounding boxes though.

	p.fbo.Bind()
	defer p.fbo.Unbind()
	p.fbo.Resize(args.Viewport.Width, args.Viewport.Height)

	// disable depth testing
	// todo: should be disabled for transparent things, not everything
	// render.DepthOutput(false)

	if err := gl.GetError(); err != nil {
		fmt.Println("uncaught GL error during pre-forward pass:", err)
	}

	objects := query.New[mesh.T]().Where(isDrawForward).Collect(scene)
	for _, drawable := range objects {
		if err := p.DrawForward(args, drawable); err != nil {
			fmt.Printf("forward draw error in object %s: %s\n", drawable.Object().Name(), err)
		}
	}

	render.DepthOutput(true)
	render.CullFace(render.CullNone)

	if err := gl.GetError(); err != nil {
		fmt.Println("uncaught GL error during forward pass:", err)
	}
}

func isDrawForward(m mesh.T) bool {
	return m.Mode() == mesh.Forward
}

func (p *ForwardPass) DrawForward(args render.Args, m mesh.T) error {
	args = args.Apply(m.Transform().World())

	mat := m.Material()

	if err := mat.Use(); err != nil {
		return fmt.Errorf("failed to assign material %s in mesh %s: %w", mat.Name(), m.Name(), err)
	}

	// set up uniforms
	mat.Mat4("model", args.Transform)
	mat.Mat4("view", args.View)
	mat.Mat4("projection", args.Projection)
	mat.Mat4("mvp", args.MVP)

	drawable := p.meshes.Fetch(m.Mesh(), mat)
	return drawable.Draw()
}
