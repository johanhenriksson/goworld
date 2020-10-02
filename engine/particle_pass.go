package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type ParticlePass struct {
}

func NewParticlePass() *ParticlePass {
	return &ParticlePass{}
}

// DrawPass executes the particle pass
func (p *ParticlePass) DrawPass(scene *Scene) {
	// kind of awkward
	scene.DrawPass(render.ParticlePass)
}

type Particle struct {
	Position vec3.T
	Velocity vec3.T
	Duration float32
}

type ParticleSystem struct {
	*Object
	Particles []Particle
	Count     int
	Chance    float32
	MinVel    vec3.T
	MaxVel    vec3.T
	MinDur    float32
	MaxDur    float32

	positions vec3.Array
	mat       *render.Material
	vao       *render.VertexArray
}

func (ps *ParticleSystem) Update(dt float32) {
	if len(ps.Particles) < ps.Count && random.Chance(ps.Chance) {
		// add particle

		p := Particle{
			Position: vec3.Zero,
			Velocity: vec3.Random(ps.MinVel, ps.MaxVel),
			Duration: random.Range(ps.MinDur, ps.MaxDur),
		}
		ps.Particles = append(ps.Particles, p)
	}

	for i := 0; i < len(ps.Particles); i++ {
		if ps.Particles[i].Duration < 0 {
			// dead
			ps.remove(i)
			i--
		}
	}

	for i, p := range ps.Particles {
		ps.Particles[i].Duration -= dt
		ps.Particles[i].Position = p.Position.Add(p.Velocity.Scaled(dt))
		ps.positions[i] = p.Position
	}

	ps.vao.Buffer("geometry", ps.positions[:len(ps.Particles)])
}

func (ps *ParticleSystem) remove(i int) {
	ps.Particles[len(ps.Particles)-1], ps.Particles[i] = ps.Particles[i], ps.Particles[len(ps.Particles)-1]
	ps.Particles = ps.Particles[:len(ps.Particles)-1]
}

func (ps *ParticleSystem) Draw(args render.DrawArgs) {
	if args.Pass != render.ParticlePass {
		return
	}

	gl.BlendFunc(gl.ONE, gl.ONE)
	gl.DepthMask(false)
	ps.mat.Use()
	ps.mat.Vec3("eye", &args.Position)
	ps.mat.Mat4("model", &args.Transform)
	ps.mat.Mat4("vp", &args.VP)
	ps.vao.Draw()
	gl.DepthMask(true)
}

func NewParticleSystem(parent *Object) *ParticleSystem {
	count := 8
	mat := assets.GetMaterial("billboard")
	ps := &ParticleSystem{
		Object: parent,
		Count:  count,
		Chance: 0.08,
		MinVel: vec3.New(-0.05, 0.4, -0.05),
		MaxVel: vec3.New(0.05, 0.6, 0.05),
		MinDur: 2,
		MaxDur: 3,

		mat:       mat,
		vao:       render.CreateVertexArray(render.Points, "geometry"),
		positions: make(vec3.Array, count),
	}
	parent.Attach(ps)

	// ps.vao.Buffer("geometry", ps.positions)
	mat.SetupVertexPointers()

	return ps
}