package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"
)

type ParticlePass struct {
}

func NewParticlePass() *ParticlePass {
	return &ParticlePass{}
}

// DrawPass executes the line pass
func (p *ParticlePass) DrawPass(scene *Scene) {
	scene.DrawPass(render.ParticlePass)
}

type ParticleSystem struct {
	*ComponentBase
	Shader    *render.ShaderProgram
	Material  *render.Material
	vao       *render.VertexArray
	vbo       *render.VertexBuffer
	Positions render.Vec3Buffer
}

func (ps *ParticleSystem) Update(dt float32) {

}

func (ps *ParticleSystem) Draw(args render.DrawArgs) {
	if args.Pass != render.ParticlePass {
		return
	}

	// update positions?
	gl.Disable(gl.CULL_FACE)
	ps.Shader.Use()
	ps.Shader.Vec3("cameraPos", &args.Position)
	ps.Shader.Mat4f("vp", &args.VP)
	ps.vao.DrawElements()
}

func NewParticleSystem(parent *Object) *ParticleSystem {
	mat := assets.GetMaterial("billboard")
	ps := &ParticleSystem{
		Positions: render.Vec3Buffer{{X: 10, Y: 10, Z: 10}},
		Material:  mat,
		Shader:    mat.ShaderProgram,
		vao:       render.CreateVertexArray(),
		vbo:       render.CreateVertexBuffer(),
	}
	ps.vao.Type = gl.POINTS

	mat.Use()
	ps.vao.Bind()
	ps.vbo.Bind()
	ps.vbo.Buffer(ps.Positions)
	mat.SetupVertexPointers()

	ps.ComponentBase = NewComponent(parent, ps)
	return ps
}
