package game

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type CollisionCheck func(*Player, vec3.T) (bool, vec3.T)

type Player struct {
	*engine.Camera

	Gravity     float32
	Speed       float32
	Airspeed    float32
	JumpForce   float32
	Friction    vec3.T
	AirFriction vec3.T
	CamHeight   vec3.T
	Flying      bool
	Grounded    bool

	collide  CollisionCheck
	position vec3.T
	velocity vec3.T
}

func NewPlayer(camera *engine.Camera, collide CollisionCheck) *Player {
	p := &Player{
		Camera:      camera,
		collide:     collide,
		Gravity:     float32(53),
		Speed:       float32(60),
		Airspeed:    float32(33),
		JumpForce:   0.25,
		Friction:    vec3.New(0.91, 1, 0.91),
		AirFriction: vec3.New(0.955, 1, 0.955),
		CamHeight:   vec3.New(0, 1.75, 0),
		Flying:      false,
	}
	p.position = camera.Position.Sub(p.CamHeight)
	return p
}

func (p *Player) Update(dt float32) {
	move := vec3.Zero
	moving := false
	if keys.Down(keys.W) && !keys.Down(keys.S) {
		move.Z += 1.0
		moving = true
	}
	if keys.Down(keys.S) && !keys.Down(keys.W) {
		move.Z -= 1.0
		moving = true
	}
	if keys.Down(keys.A) && !keys.Down(keys.D) {
		move.X -= 1.0
		moving = true
	}
	if keys.Down(keys.D) && !keys.Down(keys.A) {
		move.X += 1.0
		moving = true
	}
	if p.Flying && keys.Down(keys.Q) && !keys.Down(keys.E) {
		move.Y -= 1.0
		moving = true
	}
	if p.Flying && keys.Down(keys.E) && !keys.Down(keys.Q) {
		move.Y += 1.0
		moving = true
	}
	if keys.Pressed(keys.V) {
		p.Flying = !p.Flying
	}

	if moving {
		right := p.Camera.Right.Scaled(move.X)
		forward := p.Camera.Forward.Scaled(move.Z)
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

	if keys.Down(keys.LeftShift) {
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

	// apply movement in Y
	p.position.Y += step.Y
	step.Y = 0

	// ground collision
	if collides, point := p.collide(p, p.position); collides {
		p.position.Y = point.Y
		p.velocity.Y = 0
		p.Grounded = true
	} else {
		p.Grounded = false
	}

	// jumping
	if p.Grounded && keys.Down(keys.Space) {
		p.velocity.Y += p.JumpForce * p.Gravity
	}

	// x collision
	xstep := p.position.Add(vec3.New(step.X, 0, 0))
	// if p.world.HeightAt(xstep) > p.position.Y {
	if collides, _ := p.collide(p, xstep); collides {
		step.X = 0
	}

	// z collision
	zstep := p.position.Add(vec3.New(0, 0, step.Z))
	// if p.world.HeightAt(zstep) > p.position.Y {
	if collides, _ := p.collide(p, zstep); collides {
		step.Z = 0
	}

	// add horizontal movement
	p.position = p.position.Add(step)

	// update camera position
	p.Camera.Position = p.position.Add(p.CamHeight)
}
