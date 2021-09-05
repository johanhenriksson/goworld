package game

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type CollisionCheck func(*Player, vec3.T) (bool, vec3.T)

type Player struct {
	object.T
	Eye object.T

	Gravity     float32
	Speed       float32
	Airspeed    float32
	JumpForce   float32
	Friction    vec3.T
	AirFriction vec3.T
	CamHeight   vec3.T
	Flying      bool
	Grounded    bool

	camera    camera.T
	collide   CollisionCheck
	velocity  vec3.T
	keys      keys.State
	mouselook bool
}

func NewPlayer(position vec3.T, cam camera.T, collide CollisionCheck) *Player {
	p := &Player{
		T:           object.New("Player"),
		Eye:         object.New("Eye"),
		camera:      cam,
		collide:     collide,
		Gravity:     float32(53),
		Speed:       float32(60),
		Airspeed:    float32(33),
		JumpForce:   0.25,
		Friction:    vec3.New(0.91, 1, 0.91),
		AirFriction: vec3.New(0.955, 1, 0.955),
		CamHeight:   vec3.New(0, 1.75, 0),
		Flying:      false,
		keys:        keys.NewState(),
	}
	p.Transform().SetPosition(position)
	p.Eye.Transform().SetPosition(p.CamHeight)
	p.Adopt(p.Eye)
	p.Eye.Attach(cam)
	return p
}

func (p *Player) KeyEvent(e keys.Event) {
	p.keys.KeyEvent(e)

	// toggle flying
	if keys.Pressed(e, keys.V) {
		p.Flying = !p.Flying
	}
}

func (p *Player) Update(dt float32) {
	move := vec3.Zero
	moving := false
	if p.keys.Down(keys.W) && p.keys.Up(keys.S) {
		move.Z += 1.0
		moving = true
	}
	if p.keys.Down(keys.S) && p.keys.Up(keys.W) {
		move.Z -= 1.0
		moving = true
	}
	if p.keys.Down(keys.A) && p.keys.Up(keys.D) {
		move.X -= 1.0
		moving = true
	}
	if p.keys.Down(keys.D) && p.keys.Up(keys.A) {
		move.X += 1.0
		moving = true
	}
	if p.Flying && p.keys.Down(keys.Q) && p.keys.Up(keys.E) {
		move.Y -= 1.0
		moving = true
	}
	if p.Flying && p.keys.Down(keys.E) && p.keys.Up(keys.Q) {
		move.Y += 1.0
		moving = true
	}

	if moving {
		right := p.Eye.Transform().Right().Scaled(move.X)
		forward := p.Eye.Transform().Forward().Scaled(move.Z)
		up := vec3.New(0, move.Y, 0)

		move = right.Add(forward)
		move.Y = 0 // remove y component
		if p.Flying {
			move = move.Add(up)
		}
		move.Normalize()
	}
	if p.Grounded || p.Flying {
		move.Scale(p.Speed)
	} else {
		move.Scale(p.Airspeed)
	}

	if p.keys.Shift() {
		move.Scale(2)
	}

	// apply movement
	p.velocity = p.velocity.Add(move.Scaled(dt))

	// friction
	if p.Grounded {
		p.velocity = p.velocity.Mul(p.Friction)
	} else {
		p.velocity = p.velocity.Mul(p.AirFriction)
	}

	// gravity
	if !p.Flying {
		p.velocity.Y -= p.Gravity * dt
	} else {
		// apply Y friction while flying
		p.velocity.Y *= p.AirFriction.X
	}

	step := p.velocity.Scaled(dt)

	position := p.Transform().Position()

	// apply movement in Y
	position.Y += step.Y
	step.Y = 0

	// ground collision
	if collides, point := p.collide(p, position); collides {
		position.Y = point.Y
		p.velocity.Y = 0
		p.Grounded = true
	} else {
		p.Grounded = false
	}

	// jumping
	if p.Grounded && p.keys.Down(keys.Space) {
		p.velocity.Y += p.JumpForce * p.Gravity
	}

	// x collision
	xstep := position.Add(vec3.New(step.X, 0, 0))
	// if p.world.HeightAt(xstep) > p.position.Y {
	if collides, _ := p.collide(p, xstep); collides {
		step.X = 0
	}

	// z collision
	zstep := position.Add(vec3.New(0, 0, step.Z))
	// if p.world.HeightAt(zstep) > p.position.Y {
	if collides, _ := p.collide(p, zstep); collides {
		step.Z = 0
	}

	// add horizontal movement
	position = position.Add(step)

	// update camera position
	//p.camera.SetPosition(p.position.Add(p.CamHeight))
	p.Transform().SetPosition(position)

	p.T.Update(dt)
}

func (p *Player) MouseEvent(e mouse.Event) {
	if e.Action() == mouse.Press && e.Button() == mouse.Button1 {
		p.mouselook = true
	}
	if e.Action() == mouse.Release && e.Button() == mouse.Button1 {
		p.mouselook = false
	}
	if e.Action() == mouse.Move && p.mouselook {
		sensitivity := vec2.New(0.045, 0.04)
		delta := e.Delta().Mul(sensitivity)

		xrot := p.Eye.Transform().Rotation().X + delta.Y
		yrot := p.Eye.Transform().Rotation().Y + delta.X

		// rot := p.Eye.Transform().Rotation().XY().Sub(delta)

		// camera angle limits
		xrot = math.Clamp(xrot, -89.9, 89.9)
		yrot = math.Mod(yrot, 360)

		p.Eye.Transform().SetRotation(vec3.New(xrot, yrot, 0))
		// p.Transform().SetRotation(vec3.New(0, yrot, 0))
	}
}
