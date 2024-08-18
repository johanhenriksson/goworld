package editor

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type Player struct {
	object.Object
	Camera   *camera.Object
	Speed    float32
	Friction vec3.T

	velocity  vec3.T
	keys      keys.State
	mouselook bool
}

func NewPlayer(pool object.Pool, position vec3.T, rotation quat.T) *Player {
	p := object.Builder(object.New(pool, "Player", &Player{
		Camera: object.Builder(camera.NewObject(pool, camera.Args{
			Fov:   58.0,
			Near:  0.1,
			Far:   500,
			Clear: color.Hex("#eddaab"),
		})).
			Rotation(rotation).
			Create(),
		Speed:    float32(33),
		Friction: vec3.New(2, 2, 2),
		keys:     keys.NewState(),
	})).
		Position(position).
		Create()
	return p
}

func (p *Player) KeyEvent(e keys.Event) {
	p.keys.KeyEvent(e)
}

func (p *Player) Update(scene object.Component, dt float32) {
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
	if p.keys.Down(keys.Q) && p.keys.Up(keys.E) {
		move.Y -= 1.0
		moving = true
	}
	if p.keys.Down(keys.E) && p.keys.Up(keys.Q) {
		move.Y += 1.0
		moving = true
	}

	if moving {
		right := p.Camera.Transform().Right().Scaled(move.X)
		up := p.Camera.Transform().Up().Scaled(move.Y)
		forward := p.Camera.Transform().Forward().Scaled(move.Z)

		move = right.Add(forward).Add(up)
		move.Normalize()
	}
	move.Scale(p.Speed)

	if p.keys.Shift() {
		move.Scale(2)
	}

	// apply movement
	p.velocity = p.velocity.Add(move.Scaled(dt))

	// friction
	friction := p.velocity.Mul(p.Friction)
	p.velocity = p.velocity.Sub(friction.Scaled(dt))

	// apply movement
	step := p.velocity.Scaled(dt)
	position := p.Transform().Position()
	position = position.Add(step)
	p.Transform().SetPosition(position)

	p.Object.Update(scene, dt)
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

		eye := p.Camera.Transform().Rotation().Euler()

		xrot := eye.X + delta.Y
		yrot := eye.Y + delta.X

		// camera angle limits
		xrot = math.Clamp(xrot, -89.9, 89.9)
		yrot = math.Mod(yrot, 360)
		rot := quat.Euler(xrot, yrot, 0)

		p.Camera.Transform().SetRotation(rot)

		e.Consume()
	}
}
