package deferred

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/mesh"
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

	// render passes basically want to collect objects based on their material
	// this is probably the place where all gpu interaction should happen - even moving data to the gpu

	// collect everything that should be rendered
	// check the asset cache if its available on the gpu
	// if not, queue upload
	// else, draw

	// this way, the implementation details like whether its a VAO or VBO can be hidden from the scene components themselves
	// they only keep a description of the mesh to be rendered, along with material etc
	// the asset cache could be shared among all render passes

	objects := query.New[mesh.T]().Where(isDrawDeferred).Collect(scene)
	for _, mesh := range objects {
		if err := p.DrawDeferred(args, mesh); err != nil {
			fmt.Printf("deferred draw error in object %s: %s\n", mesh.Name(), err)
		}
	}

	if err := gl.GetError(); err != nil {
		fmt.Println("uncaught GL error during geometry pass:", err)
	}
}

func (p *GeometryPass) DrawDeferred(args render.Args, mesh mesh.T) error {
	args = args.Apply(mesh.Transform().World())

	mat := mesh.Material()

	if err := mat.Use(); err != nil {
		return fmt.Errorf("failed to assign material %s in mesh %s: %w", mat.Name(), mesh.Name(), err)
	}

	// could be updated per material
	mat.Vec3("eye", args.Position)
	mat.Mat4("view", args.View)
	mat.Mat4("projection", args.Projection)

	// must be set for each model
	mat.Mat4("model", args.Transform)
	mat.Mat4("mvp", args.MVP)

	return mesh.Vao().Draw()
}

func isDrawDeferred(m mesh.T) bool {
	return m.Mode() == mesh.Deferred
}
