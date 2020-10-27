package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type ParticleDrawable interface {
	DrawParticles(DrawArgs)
}

// ParticlePass represents the particle system draw pass
type ParticlePass struct {
	queue *DrawQueue
}

// NewParticlePass creates a new particle system draw pass
func NewParticlePass() *ParticlePass {
	return &ParticlePass{
		queue: NewDrawQueue(),
	}
}

// Type returns the render pass type identifier
func (p *ParticlePass) Type() render.Pass {
	return render.Particles
}

// Resize is called on window resize. Should update any window size-dependent buffers
func (p *ParticlePass) Resize(width, height int) {}

// DrawPass executes the particle pass
func (p *ParticlePass) Draw(scene *Scene) {
	p.queue.Clear()
	scene.Collect(p)
	for _, cmd := range p.queue.items {
		drawable := cmd.Component.(ParticleDrawable)
		drawable.DrawParticles(cmd.Args)
	}
}

func (p *ParticlePass) Visible(c Component, args DrawArgs) bool {
	_, ok := c.(ParticleDrawable)
	return ok
}

func (p *ParticlePass) Queue(c Component, args DrawArgs) {
	p.queue.Add(c, args)
}

// Particle holds data about a single particle
type Particle struct {
	Position vec3.T
	Velocity vec3.T
	Duration float32
}

// ParticleSystem holds the properties of a particle system effect
type ParticleSystem struct {
	*Transform

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

// Update the particle system
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

// Draw the particle system
func (ps *ParticleSystem) Draw(args DrawArgs) {
	if args.Pass != render.Particles {
		return
	}

	args = args.Apply(ps.Transform)

	render.Blend(true)
	render.BlendFunc(gl.ONE, gl.ONE)
	render.DepthOutput(false)

	ps.mat.Use()
	ps.mat.Vec3("eye", &args.Position)
	ps.mat.Mat4("model", &args.Transform)
	ps.mat.Mat4("vp", &args.VP)
	ps.vao.Draw()

	render.DepthOutput(true)
}

// NewParticleSystem creates a new particle system
func NewParticleSystem(position vec3.T) *ParticleSystem {
	count := 8
	mat := assets.GetMaterial("billboard")
	ps := &ParticleSystem{
		Transform: NewTransform(position, vec3.Zero, vec3.One),

		Count:  count,
		Chance: 0.08,
		MinVel: vec3.New(-0.05, 0.4, -0.05),
		MaxVel: vec3.New(0.05, 0.6, 0.05),
		MinDur: 2,
		MaxDur: 3,

		mat:       mat,
		vao:       render.CreateVertexArray(render.Points),
		positions: make(vec3.Array, count),
	}

	//mat.SetupVertexPointers()

	return ps
}
