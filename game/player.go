package game

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type CollisionCheck func(*Player, vec3.T) (bool, vec3.T)

type Player struct {
	object.T
	Eye    object.T
	Camera camera.T

	Gravity     float32
	Speed       float32
	Airspeed    float32
	JumpForce   float32
	Friction    vec3.T
	AirFriction vec3.T
	CamHeight   vec3.T
	Flying      bool
	Grounded    bool

	collide   CollisionCheck
	velocity  vec3.T
	keys      keys.State
	mouselook bool
}

func NewPlayer(position vec3.T, collide CollisionCheck) *Player {
	cam := camera.New(50.0, 0.1, 500, color.Hex("#eddaab"))
	p := object.New(&Player{
		Eye: object.Builder(object.Empty("Eye")).
			Position(vec3.New(0, 1.75, 0)).
			Attach(cam).
			Attach(light.NewPoint(light.PointArgs{
				Attenuation: light.DefaultAttenuation,
				Range:       20,
				Intensity:   2.5,
				Color:       color.White,
			})).
			Create(),
		Camera:      cam,
		collide:     collide,
		Gravity:     float32(53),
		Speed:       float32(60),
		Airspeed:    float32(33),
		JumpForce:   0.25,
		Friction:    vec3.New(3, 0, 3),
		AirFriction: vec3.New(2, 2, 2),
		CamHeight:   vec3.New(0, 1.75, 0),
		Flying:      collide == nil,
		keys:        keys.NewState(),
	})
	p.Transform().SetPosition(position)
	return p
}

func (p *Player) KeyEvent(e keys.Event) {
	p.keys.KeyEvent(e)

	// toggle flying
	if keys.Pressed(e, keys.V) && p.collide != nil {
		p.Flying = !p.Flying
	}
}

func (p *Player) Update(scene object.T, dt float32) {
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
		friction := p.velocity.Mul(p.Friction)
		p.velocity = p.velocity.Sub(friction.Scaled(dt))
		if p.velocity.Length() < 0.01 {
			p.velocity = vec3.Zero
		}
	} else {
		friction := p.velocity.Mul(p.AirFriction)
		if !p.Flying {
			friction.Y = 0
		}
		p.velocity = p.velocity.Sub(friction.Scaled(dt))
	}

	// gravity
	if !p.Flying {
		p.velocity.Y -= p.Gravity * dt
	}

	step := p.velocity.Scaled(dt)

	position := p.Transform().Position()

	// apply movement in Y
	position.Y += step.Y
	step.Y = 0

	// ground collision
	p.Grounded = false
	if p.collide != nil {
		if collides, point := p.collide(p, position); collides {
			position.Y = point.Y
			p.velocity.Y = 0
			p.Grounded = true
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
	}

	// add horizontal movement
	position = position.Add(step)

	// update camera position
	p.Transform().SetPosition(position)

	p.T.Update(scene, dt)
}

func (p *Player) MouseEvent(e mouse.Event) {
	if e.Action() == mouse.Press && e.Button() == mouse.Button2 {
		p.mouselook = true
		mouse.Lock()
		e.Consume()
	}
	if e.Action() == mouse.Release && e.Button() == mouse.Button2 {
		p.mouselook = false
		mouse.Show()
		e.Consume()
	}

	if e.Action() == mouse.Move && p.mouselook {
		sensitivity := vec2.New(0.045, 0.04)
		delta := e.Delta().Mul(sensitivity)

		eye := p.Eye.Transform().Rotation().Euler()

		xrot := eye.X - delta.Y
		yrot := eye.Y - delta.X

		// camera angle limits
		xrot = math.Clamp(xrot, -89.9, 89.9)
		yrot = math.Mod(yrot, 360)
		rot := quat.Euler(xrot, yrot, 0)

		p.Eye.Transform().SetRotation(rot)

		e.Consume()
	}
}
